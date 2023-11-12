package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func ListSessions(c *fiber.Ctx) error {
	user := c.Locals("session").(models.Session)

	var sessions []models.Session
	if err := lib.DB.Where(models.Session{
		Connection: &models.Connection{
			UserID: user.Connection.UserID,
		},
	}).Find(&sessions).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    sessions,
	})
}
