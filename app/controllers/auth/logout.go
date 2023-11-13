package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
	"net/http"
)

func Logout(c *fiber.Ctx) error {
	authCookie := c.Cookies("Authorization")

	lib.ClearAuth(c)

	lib.DB.Delete(&models.Session{
		ID: authCookie,
	})

	return c.Status(http.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    nil,
	})
}
