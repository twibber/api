package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type UserResponse struct {
	models.User
	Counts Counts `json:"counts"`
}

type Counts struct {
	Followers int `json:"followers"`
	Following int `json:"following"`

	// gorm uses int64 for counts
	Posts int64 `json:"posts"`
	Likes int64 `json:"likes"`
}

func ListUsers(c *fiber.Ctx) error {
	session := lib.GetSession(c)
	userID := ""

	if session != nil {
		userID = session.Connection.User.ID
	}

	var dbUsers []models.User
	if err := lib.DB.
		Model(&models.User{}).
		Preload("Followers").
		Preload("Following").
		Order("users.created_at DESC").
		Find(&dbUsers).Error; err != nil {
		return err
	}

	// check if the user follows each user and if each user follows the user
	for i, dbUser := range dbUsers {
		for _, follower := range dbUser.Followers {
			if follower.UserID == userID {
				dbUsers[i].YouFollow = true
				break
			}
		}

		for _, following := range dbUser.Following {
			if following.UserID == userID {
				dbUsers[i].FollowsYou = true
				break
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    dbUsers,
	})
}

func GetUserByUsername(c *fiber.Ctx) error {
	session := lib.GetSession(c)

	curUserID := ""
	if session != nil {
		curUserID = session.Connection.User.ID
	}

	username := c.Params("user")

	var user models.User
	if err := lib.DB.
		Where("username = ?", username).
		Preload("Followers").
		Preload("Following").
		First(&user).Error; err != nil {
		return err
	}

	// Check if the user follows each user and if each user follows the user
	for _, follower := range user.Followers {
		if follower.UserID == curUserID {
			user.YouFollow = true
			break
		}

		for _, following := range user.Following {
			if following.UserID == curUserID {
				user.FollowsYou = true
				break
			}
		}
	}

	// Count the number of followers, following, posts, and likes
	var counts Counts

	counts.Followers = len(user.Followers)
	counts.Following = len(user.Following)

	if err := lib.DB.Model(&models.Post{}).Where("user_id = ?", user.ID).Count(&counts.Posts).Error; err != nil {
		return err
	}

	if err := lib.DB.Model(&models.Like{}).Where("user_id = ?", user.ID).Count(&counts.Likes).Error; err != nil {
		return err
	}

	// Respond with the user data
	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data: UserResponse{
			User:   user,
			Counts: counts,
		},
	})
}
