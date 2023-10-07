package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type VerifyDTO struct {
	Email string `json:"email" validate:"required,email,max=512"`
	Code  string `json:"code" validate:"required"`
}

func Verify(c *fiber.Ctx) error {
	var body VerifyDTO
	if err := lib.ParseAndValidate(c, &body); err != nil {
		return err
	}

	var user models.User
	if err := lib.DB.Where(models.User{Email: body.Email}).First(&user).Error; err != nil {
		return err
	}

	if !lib.ValidateTOTP(user.MFA, body.Code, lib.EmailVerification) {
		return lib.NewError(fiber.StatusBadRequest, "Invalid code provided.", &lib.ErrorDetails{
			Fields: []lib.ErrorField{
				{
					Name:   "code",
					Errors: []string{"The code provided is invalid."},
				},
			},
		})
	} else {
		user.Verified = true
	}

	if err := lib.DB.Updates(&user).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
