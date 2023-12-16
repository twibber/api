package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func DeleteUser(c *fiber.Ctx) error {
	session := lib.GetSession(c)

	if session.Connection.User.Username == c.Params("user") {
		return lib.NewError(fiber.StatusBadRequest, "You cannot delete yourself.", nil)
	}

	var user models.User
	if err := lib.DB.Where("username = ?", c.Params("user")).First(&user).Error; err != nil {
		return err
	}

	if err := lib.DB.Delete(&user).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
