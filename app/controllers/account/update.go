package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type UpdateProfileDTO struct {
	DisplayName string `json:"display_name" validate:"omitempty,min=3,max=32,notblank"`
	Username    string `json:"username" validate:"omitempty,min=3,max=32,lowercase,ascii,notblank"`
}

func UpdateProfile(c *fiber.Ctx) error {
	session := lib.GetSession(c)

	var dto UpdateProfileDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	user := session.Connection.User

	if dto.Username != "" && dto.Username != user.Username {
		// check if user with username already exists
		var userExists int64
		if err := lib.DB.Table("users").Where(&models.User{
			Username: dto.Username,
		}).Count(&userExists).Error; err != nil {
			return err
		}

		if userExists > 0 {
			return lib.NewError(fiber.StatusBadRequest, "A user with that username already exists.", nil)
		}

		user.Username = dto.Username
	}

	if dto.DisplayName != "" && dto.DisplayName != user.DisplayName {
		user.DisplayName = dto.DisplayName
	}

	if err := lib.DB.Table("users").Where(&models.User{
		BaseModel: models.BaseModel{ID: session.Connection.User.ID},
	}).Updates(user).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
