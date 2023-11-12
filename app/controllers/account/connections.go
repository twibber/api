package account

import (
	"github.com/gofiber/fiber/v2"   // Fiber web framework for Go
	"github.com/twibber/api/lib"    // Shared library for common functionality
	"github.com/twibber/api/models" // Data models used in the application
)

// ListConnections handles the request to list all connections for the authenticated user.
func ListConnections(c *fiber.Ctx) error {
	// Retrieve the user session from the context.
	user := c.Locals("session").(models.Session)

	// Prepare a slice to hold the user's connections.
	var connections []models.Connection

	// Query the database for connections belonging to the user.
	if err := lib.DB.Where(models.Connection{
		UserID: user.Connection.UserID,
	}).Find(&connections).Error; err != nil {
		// In case of a query error, return the error.
		return err
	}

	// Return the list of connections with a success status.
	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    connections,
	})
}
