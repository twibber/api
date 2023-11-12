package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
	"net/http"
)

func Logout(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	// if it reaches this far there ain't a point to backing out
	lib.ClearAuth(c)

	if err := lib.DB.Delete(&session).Error; err != nil {
		return err
	}

	return c.Redirect(lib.Config.PublicURL, http.StatusTemporaryRedirect)
}
