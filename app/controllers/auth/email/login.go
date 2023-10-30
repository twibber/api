package email

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
	var dto LoginDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	if err := lib.CheckCaptcha(dto.Captcha); err != nil {
		return err
	}

	var connection models.Connection
	if err := lib.DB.Where(models.Connection{ID: models.Email.WithID(dto.Email)}).First(&connection).Error; err != nil {
		return err
	}
	
	match, err := argon2id.ComparePasswordAndHash(dto.Password, connection.Password)
	if err != nil {
		return err
	}

	if !match {
		return lib.ErrInvalidCredentials
	}

	token := lib.GenerateString(64)

	exp := 24 * time.Hour
	if dto.Remember {
		exp = 2 * 7 * 24 * time.Hour
	}

	lib.DB.Create(&models.Session{
		ID:           token,
		ConnectionID: connection.ID,
		ExpiresAt:    time.Now().Add(exp),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
