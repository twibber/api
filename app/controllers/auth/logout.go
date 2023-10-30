package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func Logout(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	// if it reaches this far there ain't a point to backing out
	c.ClearCookie("Authorization")

	if err := lib.DB.Delete(&session).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
