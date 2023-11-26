package middleware

import (
	"github.com/gofiber/fiber/v2"   // Web framework for Golang
	"github.com/twibber/api/lib"    // Contains shared configurations and utilities
	"github.com/twibber/api/models" // Data models for the application
)

// Auth middleware to handle authentication with optional verification.
func Auth(verify bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := lib.GetSession(c)

		if session == nil {
			lib.ClearAuth(c)
			return lib.ErrUnauthorised
		}

		// Proceeds if user is verified or verification isn't required
		if !verify || session.Connection.Verified {
			c.Locals("session", *session) // * as we don't need a pointer to the session
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
