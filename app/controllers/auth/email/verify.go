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

	if !lib.ValidateTOTP(session.Connection.TOTPVerify, dto.Code, lib.EmailVerification) {
		return lib.NewError(fiber.StatusBadRequest, "Invalid code provided.", &lib.ErrorDetails{
			Fields: []lib.ErrorField{
				{
					Name:   "code",
					Errors: []string{"The code provided is invalid."},
				},
			},
		})
	} else {
		session.Connection.Verified = true
	}

	if err := lib.DB.Updates(&session.Connection).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
