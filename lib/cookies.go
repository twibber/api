package lib

import (
	"github.com/gofiber/fiber/v2"
	cfg "github.com/twibber/api/config"
	"time"
) // Fiber web framework for Go

// ClearAuth removes the "Authorization" cookie from the client, effectively logging the user out.
func ClearAuth(c *fiber.Ctx) {
	// Sets the "Authorization" cookie to an empty value and expires it immediately.
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",   // Name of the cookie to clear
		Value:    "",                // Clears the value of the cookie
		Path:     "/",               // Path for which the cookie is valid
		Domain:   cfg.Config.Domain, // Domain for which the cookie is valid
		MaxAge:   0,                 // Sets the cookie to expire immediately
		HTTPOnly: true,              // Prevents JavaScript from accessing the cookie
		SameSite: "lax",             // Lax same-site policy to allow sending the cookie along with cross-site requests
	})
}

// SetAuth sets the "Authorization" cookie to the provided token and expires it after the provided duration.
func SetAuth(c *fiber.Ctx, token string, exp time.Duration) {
	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    token,                   // Sets the auth token as the value of the cookie
		Path:     "/",                     // Allows all paths to access the cookie
		Domain:   "." + cfg.Config.Domain, // Make the cookie a wildcard  to access the cookie from all subdomains of the domain
		MaxAge:   int(exp.Seconds()),      // Sets the cookie to expire after the provided duration
		HTTPOnly: true,                    // Prevents JavaScript from accessing the cookie
		SameSite: "lax",                   // Lax same-site policy to allow sending the cookie along with cross-site requests on the same domain
	})
}
