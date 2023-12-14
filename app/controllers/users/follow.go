package users

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
	"gorm.io/gorm"
)

func FollowUser(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var user models.User
	if err := lib.DB.Table("users").Where(models.User{
		BaseModel: models.BaseModel{ID: c.Params("user")},
	}).First(&user).Error; err != nil {
		return err
	}

	if user.ID == session.Connection.User.ID {
		return lib.NewError(fiber.StatusBadRequest, "You cannot follow yourself.", nil)
	}

	// ensure that the user does not already follow the user
	var followExists int64
	if err := lib.DB.Table("follows").Where(&models.Follow{
		UserID:     session.Connection.User.ID,
		FollowedID: user.ID,
	}).Count(&followExists).Error; err != nil {
		return err
	}

	if followExists > 0 {
		return lib.NewError(fiber.StatusBadRequest, "You are already following this user.", nil)
	}

	if err := lib.DB.Table("follows").Create(&models.Follow{
		UserID:     session.Connection.User.ID,
		FollowedID: user.ID,
	}).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}

func UnfollowUser(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var user models.User
	if err := lib.DB.Table("users").Where(models.User{
		BaseModel: models.BaseModel{ID: c.Params("user")},
	}).First(&user).Error; err != nil {
		return err
	}

	var follow models.Follow
	if err := lib.DB.Table("follows").Where(&models.Follow{
		UserID:     session.Connection.User.ID,
		FollowedID: user.ID,
	}).First(&follow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return lib.NewError(fiber.StatusBadRequest, "You are not following this user.", nil)
		} else {
			return err
		}
	}

	if err := lib.DB.Table("follows").Delete(&follow).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}

func GetFollowers(c *fiber.Ctx) error {
	var user models.User
	if err := lib.DB.Table("users").Where(models.User{
		BaseModel: models.BaseModel{ID: c.Params("user")},
	}).First(&user).Error; err != nil {
		return err
	}

	var followers []models.Follow
	if err := lib.DB.Table("follows").Where(&models.Follow{
		FollowedID: user.ID,
	}).Find(&followers).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    followers,
	})
}

func GetFollowing(c *fiber.Ctx) error {
	var user models.User
	if err := lib.DB.Table("users").Where(models.User{
		BaseModel: models.BaseModel{ID: c.Params("user")},
	}).First(&user).Error; err != nil {
		return err
	}

	var following []models.Follow
	if err := lib.DB.Table("follows").Where(&models.Follow{
		UserID: user.ID,
	}).Find(&following).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    following,
	})
}
