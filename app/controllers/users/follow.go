package users

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
	"gorm.io/gorm"
)

func FollowUser(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var user models.User
	if err := lib.DB.First(&user, c.Params("user")).Error; err != nil {
		return err
	}

	if user.ID == session.Connection.User.ID {
		return lib.NewError(fiber.StatusBadRequest, "You cannot follow yourself.", nil)
	}

	if err := lib.DB.Create(&models.Follow{
		ID:         utils.UUIDv4(),
		UserID:     session.Connection.User.ID,
		FollowedID: user.ID,
		Timestamps: lib.NewDBTime(),
	}).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}

func UnfollowUser(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var user models.User
	if err := lib.DB.First(&user, c.Params("user")).Error; err != nil {
		return err
	}

	var follow models.Follow
	if err := lib.DB.Where(&models.Follow{
		UserID:     session.Connection.User.ID,
		FollowedID: user.ID,
	}).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.NewError(fiber.StatusBadRequest, "You are not following this user.", nil)
		} else {
			return err
		}
	}

	if err := lib.DB.Delete(&follow).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
