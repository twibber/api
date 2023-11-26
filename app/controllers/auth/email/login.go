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

	tx := lib.DB.Begin()

	var connection models.Connection
	if err := tx.Where(models.Connection{ID: models.EmailType.WithID(dto.Email)}).First(&connection).Error; err != nil {
		tx.Rollback()
		return err
	}

	match, err := argon2id.ComparePasswordAndHash(dto.Password, connection.Password)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !match {
		tx.Rollback()
		return lib.ErrInvalidCredentials
	}

	token := lib.GenerateString(64)
	exp := 24 * time.Hour

	if err := tx.Create(&models.Session{
		ID:           token,
		ConnectionID: connection.ID,
		Info: models.SessionInfo{
			IPAddress: c.IP(),
			UserAgent: c.Get("User-Agent"),
		},
		ExpiresAt:  time.Now().Add(exp),
		Timestamps: lib.NewDBTime(),
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

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
