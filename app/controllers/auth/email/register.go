package email

import (
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"runtime"

	"github.com/twibber/api/lib"
	"github.com/twibber/api/mailer"
	"github.com/twibber/api/models"

	log "github.com/sirupsen/logrus"

	"net/http"
	"time"
)

type RegisterDTO struct {
	Username string `json:"username" validate:"required,max=32"`
	Email    string `json:"email"     validate:"required,email,max=512"`
	Password string `json:"password"  validate:"required,min=8"`
	Captcha  string `json:"captcha"   validate:""`
}

func Register(c *fiber.Ctx) error {
	var dto RegisterDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	if err := lib.CheckCaptcha(dto.Captcha); err != nil {
		return err
	}

	var count int64
	if err := lib.DB.Model(models.User{}).Where(models.User{Email: dto.Email}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return lib.ErrEmailExists
	}

	hashedPassword, err := argon2id.CreateHash(dto.Password, &argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  16,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  32,
		KeyLength:   128,
	})

	// check if username already exists
	var usernameCount int64
	if err := lib.DB.Model(models.User{}).Where(models.User{Username: dto.Username}).Count(&usernameCount).Error; err != nil {
		return err
	}

	if usernameCount > 0 {
		return lib.NewError(fiber.StatusConflict, "Username already exists.", &lib.ErrorDetails{
			Fields: []lib.ErrorField{
				{
					Name:   "username",
					Errors: []string{"The username provided is already in use."},
				},
			},
		})
	}

	token := lib.GenerateString(64)

	totpCode, err := lib.GenerateSecureRandomBase32(32)
	if err != nil {
		return err
	}

	code, err := lib.GenerateTOTP(totpCode, lib.EmailVerification)
	if err != nil {
		return err
	}

	tx := lib.DB.Begin()

	user := models.User{
		ID:          utils.UUIDv4(),
		Username:    dto.Username,
		DisplayName: dto.Username,
		Email:       dto.Email,
		Suspended:   false,
		Timestamps:  lib.NewDBTime(),
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}

	exp := 24 * time.Hour
	if err := tx.Create(&models.Connection{
		ID:         models.EmailType.WithID(dto.Email),
		UserID:     user.ID,
		Password:   hashedPassword,
		TOTPVerify: totpCode,
		Verified:   false,
		Sessions: []models.Session{
			{
				ID: token,
				Info: models.SessionInfo{
					IPAddress: c.IP(),
					UserAgent: c.Get("User-Agent"),
				},
				ExpiresAt:  time.Now().Add(exp),
				Timestamps: lib.NewDBTime(),
			},
		},
		Timestamps: lib.NewDBTime(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// concurrently send verification email to user
	go func() {
		err := mailer.VerifyDTO{
			Defaults: mailer.Defaults{
				Email: user.Email,
				Name:  user.Username,
			},
			Code: code,
		}.Send()
		if err != nil {
			log.WithError(err).Error("new verification code could not be sent")
		}
	}()

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		Domain:   "." + lib.Config.Domain, // adds a dot to the domain to allow subdomains
		MaxAge:   int(exp.Seconds()),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return c.Status(http.StatusOK).JSON(lib.BlankSuccess)
}
