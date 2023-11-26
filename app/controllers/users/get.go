package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func ListUsers(c *fiber.Ctx) error {
	var users []models.User
	if err := lib.DB.Find(&users).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    users,
	})
}

type UserResponse struct {
	models.User
	Counts     Counts `json:"counts"`
	Following  bool   `json:"following"`
	FollowsYou bool   `json:"follows_you"`
}

type Counts struct {
	Followers int64 `json:"followers"`
	Following int64 `json:"following"`
	Posts     int64 `json:"posts"`
	Likes     int64 `json:"likes"`
}

func GetUser(c *fiber.Ctx) error {
	// get current user
	var session = c.Locals("session").(models.Session)

	// create response
	resp := UserResponse{}

	// get user
	if err := lib.DB.Where(&models.User{
		Username: c.Params("user"),
	}).First(&resp.User).Error; err != nil {
		return err
	}

	// count followers, following, posts, likes
	var counts Counts
	if err := lib.DB.
		Model(&models.Follow{}).
		Where(&models.Follow{
			FollowedID: resp.User.ID,
		}).Count(&counts.Followers).Error; err != nil {
		return err
	}

	if err := lib.DB.
		Model(&models.Follow{}).
		Where(&models.Follow{
			UserID: resp.User.ID,
		}).Count(&counts.Following).Error; err != nil {
		return err
	}

	if err := lib.DB.
		Model(&models.Post{}).
		Where(&models.Post{
			UserID: resp.User.ID,
		}).Count(&counts.Posts).Error; err != nil {
		return err
	}

	if err := lib.DB.
		Model(&models.Like{}).
		Where(&models.Like{
			UserID: resp.User.ID,
		}).Count(&counts.Likes).Error; err != nil {
		return err
	}

	// check if the user is following you
	if err := lib.DB.
		Model(&models.Follow{}).
		Where(&models.Follow{
			UserID:     session.Connection.User.ID,
			FollowedID: resp.User.ID,
		}).First(&resp.Following).Error; err != nil {
		resp.Following = false
	}

	// check if you are following the user
	if err := lib.DB.
		Model(&models.Follow{}).
		Where(&models.Follow{
			UserID:     resp.User.ID,
			FollowedID: session.Connection.User.ID,
		}).First(&resp.FollowsYou).Error; err != nil {
		resp.FollowsYou = false
	}

	// set counts
	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    resp,
	})
}
