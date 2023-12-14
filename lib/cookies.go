package lib

import (
	"github.com/gofiber/fiber/v2"
	"time"
) // Fiber web framework for Go

// ClearAuth removes the "Authorization" cookie from the client, effectively logging the user out.
func ClearAuth(c *fiber.Ctx) {
	// Sets the "Authorization" cookie to an empty value and expires it immediately.
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization", // Name of the cookie to clear
		Value:    "",              // Clears the value of the cookie
		Path:     "/",             // Path for which the cookie is valid
		Domain:   Config.Domain,   // Domain for which the cookie is valid
		MaxAge:   0,               // Sets the cookie to expire immediately
		HTTPOnly: true,            // Prevents JavaScript from accessing the cookie
		SameSite: "lax",           // Lax same-site policy to allow sending the cookie along with cross-site requests
	})
}

func SetAuth(c *fiber.Ctx, token string, exp time.Duration) {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,
		Path:     "/",
		Domain:   "." + Config.Domain, // Allows subdomains to access the cookie
		MaxAge:   int(exp.Seconds()),
		HTTPOnly: true,
		SameSite: "lax",
	})
}
