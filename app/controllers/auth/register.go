package auth

import (
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/mailer"
	"github.com/twibber/api/models"
	"net/http"
	"time"
)

type RegisterDTO struct {
	Username string `json:"firstname" validate:"required,max=32"`
	Email    string `json:"email"     validate:"required,email,max=512"`
	Password string `json:"password"  validate:"required,min=8"`
	Remember bool   `json:"remember"  validate:""`
	Captcha  string `json:"captcha"   validate:""`
}

func Register(c *fiber.Ctx) error {
	var body RegisterDTO
	if err := lib.ParseAndValidate(c, &body); err != nil {
		return err
	}

	if err := lib.CheckCaptcha(body.Captcha); err != nil {
		return err
	}

	var count int64
	if err := lib.DB.Model(models.User{}).Where(models.User{Email: body.Email}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return lib.ErrEmailExists
	}

	hashedPassword, err := argon2id.CreateHash(body.Password, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	secret := lib.GenerateString(32)
	token := lib.GenerateString(64)

	user := models.User{
		ID:        utils.UUIDv4(),
		Username:  body.Username,
		Email:     body.Email,
		Verified:  false,
		MFA:       secret,
		Suspended: false,
		Connections: []models.Connection{
			{
				Type:     models.Email,
				ID:       body.Email,
				Password: hashedPassword,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := lib.DB.Create(&user).Error; err != nil {
		return err
	}

	exp := 24 * time.Hour
	if body.Remember {
		exp = 2 * 7 * 24 * time.Hour
	}

	if err := lib.DB.Create(&models.Session{
		ID:        token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(exp),
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
