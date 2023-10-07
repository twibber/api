package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/mailer"
	"github.com/twibber/api/models"
	"strings"
)

type ResendDTO struct {
	Email string `json:"email" validate:"required,email,max=512"`
}

func ResendCode(c *fiber.Ctx) error {
	var body ResendDTO
	if err := lib.ParseAndValidate(c, &body); err != nil {
		return err
	}

	code := strings.Split(utils.UUIDv4(), "-")[0]

	var user models.User
	if err := lib.DB.Model(models.User{}).Where(models.User{Email: body.Email}).First(&user).Error; err != nil {
		return err
	}

	if user.Email == "" {
		return lib.NewError(fiber.StatusNotFound, "No user found with that email address.", nil, "NOT_FOUND")
	}

	if user.Verified {
		return lib.NewError(fiber.StatusForbidden, "You have already verified your email.", nil, "ALREADY_VERIFIED")
	}

	// create verification code
	code, err := lib.GenerateTOTP(user.MFA, lib.EmailVerification)
	if err != nil {
		return err
	}

	// concurrently send verification email to user
	go func() {
		err := mailer.VerifyDTO{
			Defaults: mailer.Defaults{
				Email: body.Email,
				Name:  user.Username,
			},
			Code: code,
		}.Send()
		if err != nil {
			log.WithError(err).Error("new verification code could not be sent")
		}
	}()

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
