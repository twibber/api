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
		Username: c.Params("user"),
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
		Username: c.Params("user"),
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

func GetFollowersByUsername(c *fiber.Ctx) error {
	username := c.Params("user")
	var user models.User
	if err := lib.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	curUserID := ""
	session := lib.GetSession(c)
	if session != nil {
		curUserID = session.Connection.User.ID
	}

	var followers []models.Follow
	if err := lib.DB.
		Where(&models.Follow{
			FollowedID: user.ID,
		}).
		Preload("User").
		Preload("User.Followers").
		Preload("User.Following").
		Find(&followers).Error; err != nil {
		return err
	}

	var followersUsers []models.User
	for _, follower := range followers {
		follower.User.YouFollow = false
		follower.User.FollowsYou = false

		for _, followerFollower := range follower.User.Followers {
			if followerFollower.UserID == curUserID {
				follower.User.YouFollow = true
				break
			}
		}

		for _, followerFollowing := range follower.User.Following {
			if followerFollowing.FollowedID == curUserID {
				follower.User.FollowsYou = true
				break
			}
		}

		followersUsers = append(followersUsers, follower.User)
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    followersUsers,
	})
}

func GetFollowingByUsername(c *fiber.Ctx) error {
	username := c.Params("user")
	var user models.User
	if err := lib.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return err
	}

	curUserID := ""
	session := lib.GetSession(c)
	if session != nil {
		curUserID = session.Connection.User.ID
	}

	var following []models.Follow
	if err := lib.DB.
		Where(&models.Follow{
			UserID: user.ID,
		}).
		Preload("Followed").
		Preload("Followed.Followers").
		Preload("Followed.Following").
		Find(&following).Error; err != nil {
		return err
	}

	var followingUsers []models.User
	for _, followed := range following {
		followed.Followed.YouFollow = false
		followed.Followed.FollowsYou = false

		for _, followedFollower := range followed.Followed.Followers {
			if followedFollower.UserID == curUserID {
				followed.Followed.YouFollow = true
				break
			}
		}

		for _, followedFollowing := range followed.Followed.Following {
			if followedFollowing.FollowedID == curUserID {
				followed.Followed.FollowsYou = true
				break
			}
		}

		followingUsers = append(followingUsers, followed.Followed)
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    followingUsers,
	})
}
