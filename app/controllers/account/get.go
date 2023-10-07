package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func GetAccount(c *fiber.Ctx) error {
	user := c.Locals("user").(models.User)

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    user,
	})
}
