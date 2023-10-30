package mw

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func Auth(verify bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		authCookie := c.Cookies("Authorization")

		if authHeader == "" && authCookie == "" {
			log.Debug("No auth header or cookie")
			return lib.ErrUnauthorised
		}

		var authToken string

		if authHeader != "" {
			authParsed := strings.Split(authHeader, " ")
			if len(authParsed) < 2 {
				log.Debug("Invalid auth header")
				return lib.ErrUnauthorised
			}

			authToken = authParsed[1]
		} else {
			if authCookie != "" {
				authToken = authCookie
			} else {
				log.Debug("No auth header or cookie")
				return lib.ErrUnauthorised
			}
		}

		var session models.Session
		if err := lib.DB.Where(models.Session{ID: authToken}).Preload("Connection").Preload("Connection.User").First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Debug("Session not found")
				return lib.ErrUnauthorised
			} else {
				return err
			}
		}

		if session.ExpiresAt.Before(time.Now()) {
			return lib.ErrUnauthorised
		}

		// check if user is verified if the action requires a verified user
		if !verify || session.Connection.Verified {
			c.Locals("session", session)
			return c.Next()
		} else {
			return lib.NewError(fiber.StatusForbidden, "You must verify your email address before performing this action.", nil, "UNVERIFIED")
		}
	}
}

func AdminCheck(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)
	if session.Connection.User.Level != models.Admin && lib.Config.Debug == false {
		return lib.ErrForbidden
	}

	return c.Next()
}
