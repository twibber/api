package account

import (
	"github.com/gofiber/fiber/v2"   // Fiber web framework for Go
	"github.com/twibber/api/lib"    // Shared library for common functionality
	"github.com/twibber/api/models" // Data models used in the application
)

// GetAccount handles the request to retrieve the currently authenticated user's account details.
func GetAccount(c *fiber.Ctx) error {
	// Extract the user session from the context, which contains the user's account details.
	user := c.Locals("session").(models.Session)

	// Respond with the user's session data encapsulated in a standard response structure.
	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    user,
	})
}
