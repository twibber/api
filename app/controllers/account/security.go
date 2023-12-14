package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
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
