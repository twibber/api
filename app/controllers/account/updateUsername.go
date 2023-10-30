package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type UpdateUsernameDTO struct {
	Username string `json:"username" validate:"required,min=3,max=32,alphanumunicode"`
}

func UpdateUsername(c *fiber.Ctx) error {
	user := c.Locals("session").(models.Session)

	var dto UpdateUsernameDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	if err := lib.DB.Model(&user.Connection.User).Update("username", dto.Username).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success:    true,
		ObjectName: "user",
		Data:       user,
	})
}
