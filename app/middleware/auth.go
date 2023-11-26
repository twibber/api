package middleware

import (
	"errors"
	"github.com/sirupsen/logrus" // Structured logging library
	"gorm.io/gorm"               // ORM library for Golang
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"   // Web framework for Golang
	"github.com/twibber/api/lib"    // Contains shared configurations and utilities
	"github.com/twibber/api/models" // Data models for the application
)

// Auth middleware to handle authentication with optional verification.
func Auth(verify bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extracts auth token from header or cookie
		authHeader := c.Get("Authorization")
		authCookie := c.Cookies("Authorization")

		// Rejects request if no auth details are provided
		if authHeader == "" && authCookie == "" {
			logrus.Debug("No auth header or cookie")
			lib.ClearAuth(c)
			return lib.ErrUnauthorised
		}

		var authToken string

		// Parses auth token from header or uses cookie
		if authHeader != "" {
			authParsed := strings.Split(authHeader, " ")
			if len(authParsed) < 2 {
				logrus.Debug("Invalid auth header")
				lib.ClearAuth(c)
				return lib.ErrUnauthorised
			}
			authToken = authParsed[1]
		} else if authCookie != "" {
			authToken = authCookie
		}

		// Fetches session from database using auth token
		var session models.Session
		if err := lib.DB.Where(models.Session{ID: authToken}).Preload("Connection").Preload("Connection.User").First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logrus.Debug("Session not found")
				lib.ClearAuth(c)
				return lib.ErrUnauthorised
			}
			return err
		}

		// Rejects expired sessions
		if session.ExpiresAt.Before(time.Now()) {
			lib.DB.Delete(&session)
			lib.ClearAuth(c)
			return lib.ErrUnauthorised
		}

		// Proceeds if user is verified or verification isn't required
		if !verify || session.Connection.Verified {
			c.Locals("session", session)
			return c.Next()
		}

		// Rejects unverified users for actions requiring verification
		return lib.NewError(fiber.StatusForbidden, "You must verify your email address before performing this action.", nil, "UNVERIFIED")
	}
}

// AdminCheck middleware restricts access to admin-only functionality.
func AdminCheck(c *fiber.Ctx) error {
	// Retrieves session from context
	session := c.Locals("session").(models.Session)

	// Blocks non-admin users unless in debug mode
	if !session.Connection.User.Admin && !lib.Config.Debug {
		return lib.ErrForbidden
	}

	// Allows request to proceed for admins or in debug mode
	return c.Next()
}
