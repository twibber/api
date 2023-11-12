package lib

import "github.com/gofiber/fiber/v2"

func ClearAuth(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "",
		Path:     "/",
		Domain:   Config.Domain,
		MaxAge:   0,
		HTTPOnly: true,
		SameSite: "lax",
	})
}
