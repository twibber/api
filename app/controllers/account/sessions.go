package account

import (
	"github.com/gofiber/fiber/v2"   // Fiber web framework for Go
	"github.com/twibber/api/lib"    // Shared library for common functionality
	"github.com/twibber/api/models" // Data models used in the application
)

// ListSessions handles the request to retrieve all sessions associated with the authenticated user.
func ListSessions(c *fiber.Ctx) error {
	// Extract the user session from the context.
	user := c.Locals("session").(models.Session)

	// Prepare a slice to hold the user's sessions.
	var sessions []models.Session

	// Query the database for sessions associated with the user.
	if err := lib.DB.Where(models.Session{
		Connection: &models.Connection{
			UserID: user.Connection.UserID,
		},
	}).Find(&sessions).Error; err != nil {
		// If there is a query error, return it.
		return err
	}

	// Respond with the list of sessions in a standard response structure.
	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    sessions,
	})
}
