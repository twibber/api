package mw

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func Auth(verify bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		authCookie := c.Cookies("Authorization")

		if authHeader == "" && authCookie == "" {
			log.Info("no auth header or cookie")
			return lib.ErrUnauthorised
		}

		var authToken string

		if authHeader != "" {
			authParsed := strings.Split(authHeader, " ")
			if len(authParsed) < 2 {
				log.Info("invalid auth header")
				return lib.ErrUnauthorised
			}

			authToken = authParsed[1]
		} else {
			if authCookie != "" {
				authToken = authCookie
			} else {
				log.Info("2. no auth header or cookie")
				return lib.ErrUnauthorised
			}
		}

		var connection models.Connection
		if err := lib.DB.Where(models.Connection{ID: authToken}).Preload("User").First(&connection).Error; err != nil {
			return err
		}

		// check if user is verified if the action requires a verified user
		if !verify || connection.User.Verified {
			c.Locals("user", connection.User)
			c.Locals("userID", connection.UserID)
			return c.Next()
		} else {
			return lib.NewError(fiber.StatusForbidden, "You must verify your email before performing this action.", nil, "UNVERIFIED")
		}
	}
}
