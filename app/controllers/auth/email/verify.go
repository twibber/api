package email

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type VerifyDTO struct {
	Code string `json:"code" validate:"required"`
}

func Verify(c *fiber.Ctx) error {
	var session = c.Locals("session").(models.Session)

	var dto VerifyDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	var connection models.Connection
	if err := lib.DB.Where(models.Connection{ID: models.EmailType.WithID(session.Connection.User.Email)}).First(&connection).Error; err != nil {
		return err
	}

	if !lib.ValidateTOTP(connection.TOTPVerify, dto.Code, lib.EmailVerification) {
		return lib.NewError(fiber.StatusBadRequest, "Invalid code provided.", &lib.ErrorDetails{
			Fields: []lib.ErrorField{
				{
					Name:   "code",
					Errors: []string{"The code provided is invalid."},
				},
			},
		})
	} else {
		connection.Verified = true
	}

	if err := lib.DB.Updates(&connection).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
