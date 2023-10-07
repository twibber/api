package auth

import (
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"

	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type LoginDTO struct {
	Email    string `json:"email"     validate:"required,email,max=512"`
	Password string `json:"password"  validate:"required,min=8"`
	Remember bool   `json:"remember"  validate:""`
	Captcha  string `json:"captcha"   validate:""`
}

func Login(c *fiber.Ctx) error {
	var body LoginDTO
	if err := lib.ParseAndValidate(c, &body); err != nil {
		return err
	}

	if err := lib.CheckCaptcha(body.Captcha); err != nil {
		return err
	}

	var connection models.Connection
	if err := lib.DB.Where(models.Connection{Type: models.Email, ID: body.Email}).First(&connection).Error; err != nil {
		return err
	}

	match, err := argon2id.ComparePasswordAndHash(body.Password, connection.Password)
	if err != nil {
		return err
	}

	if !match {
		return lib.ErrInvalidCredentials
	}

	token := lib.GenerateString(64)

	exp := 24 * time.Hour
	if body.Remember {
		exp = 2 * 7 * 24 * time.Hour
	}

	lib.DB.Create(&models.Session{
		ID:        token,
		UserID:    connection.UserID,
		ExpiresAt: time.Now().Add(exp),
	})

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
