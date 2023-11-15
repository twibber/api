package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func Logout(c *fiber.Ctx) error {
	authCookie := c.Cookies("Authorization")

	lib.ClearAuth(c)

	lib.DB.Delete(&models.Session{
		ID: authCookie,
	})

	return c.Redirect(lib.Config.PublicURL, fiber.StatusTemporaryRedirect)
}
