package email

import (
	"github.com/alexedwards/argon2id"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"

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
	Remember bool   `json:"remember"  validate:""`
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

	hashedPassword, err := argon2id.CreateHash(dto.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	secret := lib.GenerateString(32)
	token := lib.GenerateString(64)

	user := models.User{
		ID:        utils.UUIDv4(),
		Username:  dto.Username,
		Email:     dto.Email,
		MFA:       secret,
		Suspended: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := lib.DB.Create(&user).Error; err != nil {
		return err
	}

	exp := 24 * time.Hour
	if dto.Remember {
		exp = 2 * 7 * 24 * time.Hour
	}

	if err := lib.DB.Create(&models.Connection{
		ID:        models.Email.WithID(dto.Email),
		UserID:    user.ID,
		Password:  hashedPassword,
		Verified:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Sessions: []models.Session{
			{
				ID:        token,
				ExpiresAt: time.Now().Add(exp),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}).Error; err != nil {
		return err
	}

	// create verification code
	code, err := lib.GenerateTOTP(secret, lib.EmailVerification)
	if err != nil {
		return err
	}

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
		Domain:   lib.Config.Domain,
		MaxAge:   int(exp.Seconds()),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return c.Status(http.StatusNoContent).JSON(lib.BlankSuccess)
}
