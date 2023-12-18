package account

import (
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
	"strings"
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
			UserID: user.Connection.User.ID,
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

// GetSession handles the request to retrieve a single session associated with the authenticated user.
func GetSession(c *fiber.Ctx) error {
	// Extract the user session from the context.
	user := c.Locals("session").(models.Session)

	// Extract the session ID from the request parameters.
	sessionID := c.Params("session")

	// Prepare a variable to hold the session.
	var session models.Session

	// Query the database for the session.
	if err := lib.DB.Where(models.Session{
		BaseModel: models.BaseModel{
			ID: sessionID,
		},
		Connection: &models.Connection{
			UserID: user.Connection.UserID,
		},
	}).First(&session).Error; err != nil {
		// If there is a query error, return it.
		return err
	}

	// Respond with the session in a standard response structure.
	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    session,
	})
}

// DeleteSession handles the request to delete a single session associated with the authenticated user.
func DeleteSession(c *fiber.Ctx) error {
	// Extract the user session from the context.
	userSession := c.Locals("session").(models.Session)

	// Extract the session ID from the request parameters.
	sessionID := c.Params("session")

	if sessionID == userSession.ID {
		return lib.NewError(fiber.StatusBadRequest, "Cannot delete the current session", nil)
	}

	// Prepare a variable to hold the session.
	var session models.Session

	// Query the database for the session.
	if err := lib.DB.Where(models.Session{
		BaseModel: models.BaseModel{
			ID: sessionID,
		},
		Connection: &models.Connection{
			UserID: userSession.Connection.UserID,
		},
	}).First(&session).Error; err != nil {
		// If there is a query error, return it.
		return err
	}

	// Delete the session from the database.
	if err := lib.DB.Delete(&session).Error; err != nil {
		// If there is a delete error, return it.
		return err
	}

	// Respond with a success status.
	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
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

// GetConnection handles the request to retrieve a single connection for the authenticated user.
func GetConnection(c *fiber.Ctx) error {
	// Retrieve the user session from the context.
	user := c.Locals("session").(models.Session)

	// Retrieve the connection ID from the request parameters.
	connID := c.Params("connection")

	// Prepare a variable to hold the connection.
	var connection models.Connection

	// Query the database for the connection.
	if err := lib.DB.Where(models.Connection{
		BaseModel: models.BaseModel{
			ID: connID,
		},
		UserID: user.Connection.UserID,
	}).First(&connection).Error; err != nil {
		// If there is a query error, return it.
		return err
	}

	// Return the connection with a success status.
	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    connection,
	})
}

type UpdateConnectionPasswordDTO struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UpdateConnectionPassword(c *fiber.Ctx) error {
	// Retrieve the user session from the context.
	user := c.Locals("session").(models.Session)

	connID := c.Params("connection")

	provider := strings.Split(connID, ":")[0]
	if provider != string(models.ProviderEmailType) {
		return lib.NewError(fiber.StatusBadRequest, "Invalid provider", nil)
	}

	// Prepare a DTO to hold the request body.
	var dto UpdateConnectionPasswordDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	// Prepare a variable to hold the connection.
	var connection models.Connection
	if err := lib.DB.Where(models.Connection{
		BaseModel: models.BaseModel{
			ID: connID,
		},
		UserID: user.Connection.UserID,
	}).First(&connection).Error; err != nil {
		return err
	}

	// check if the old password is correct
	match, err := argon2id.ComparePasswordAndHash(dto.OldPassword, connection.Password)
	if err != nil {
		return err
	}

	if !match {
		return lib.NewError(fiber.StatusBadRequest, "Incorrect password", &lib.ErrorDetails{
			Fields: []lib.ErrorField{
				{
					Name:   "old_password",
					Errors: []string{"The password is incorrect"},
				},
			},
		})
	}

	// Hash the new password
	passwordHash, err := argon2id.CreateHash(dto.NewPassword, &lib.ArgonConfig)

	if err := lib.DB.Model(&connection).Update("password", passwordHash).Error; err != nil {
		return err
	}

	// Return a success response.
	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
