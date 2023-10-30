package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func GetAccount(c *fiber.Ctx) error {
	user := c.Locals("session").(models.Session)

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success:    true,
		ObjectName: "user",
		Data:       user,
	})
}
