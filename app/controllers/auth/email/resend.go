package email

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
	var dto ResendDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	code := strings.Split(utils.UUIDv4(), "-")[0]

	var connection models.Connection
	if err := lib.DB.Where(models.Connection{ID: models.EmailType.WithID(dto.Email)}).Preload("User").First(&connection).Error; err != nil {
		return err
	}

	if dto.Email == "" {
		return lib.NewError(fiber.StatusNotFound, "No user found with that email address.", nil, "NOT_FOUND")
	}

	if connection.Verified {
		return lib.NewError(fiber.StatusForbidden, "You have already verified your email.", nil, "ALREADY_VERIFIED")
	}

	// create verification code
	code, err := lib.GenerateTOTP(connection.TOTPVerify, lib.EmailVerification)
	if err != nil {
		return err
	}

	// concurrently send verification email to user
	go func() {
		err := mailer.VerifyDTO{
			Defaults: mailer.Defaults{
				Email: dto.Email,
				Name:  connection.User.Username,
			},
			Code: code,
		}.Send()
		if err != nil {
			log.WithError(err).Error("new verification code could not be sent")
		}
	}()

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
