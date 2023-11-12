package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func ListConnections(c *fiber.Ctx) error {
	user := c.Locals("session").(models.Session)

	var connections []models.Connection
	if err := lib.DB.Where(models.Connection{
		UserID: user.Connection.UserID,
	}).Find(&connections).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    connections,
	})
}
