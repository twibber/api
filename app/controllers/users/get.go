package users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type UserQueryResult struct {
	models.User

	CountsFollowers int64 `gorm:"column:counts_followers"`
	CountsFollowing int64 `gorm:"column:counts_following"`
	CountsPosts     int64 `gorm:"column:counts_posts"`
	CountsLikes     int64 `gorm:"column:counts_likes"`
}

type UserResponse struct {
	models.User
	Counts Counts `json:"counts"`
}

type Counts struct {
	Followers int64 `json:"followers"`
	Following int64 `json:"following"`
	Posts     int64 `json:"posts"`
	Likes     int64 `json:"likes"`
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

	// get user by username
	var dbUser UserQueryResult
	if err := lib.DB.
		Table("users").
		Select("users.*, "+
			"COALESCE((SELECT COUNT(*) FROM follows WHERE followed_id = users.id), 0) as counts_followers, "+
			"COALESCE((SELECT COUNT(*) FROM follows WHERE user_id = users.id), 0) as counts_following, "+
			"COALESCE((SELECT COUNT(*) FROM posts WHERE user_id = users.id), 0) as counts_posts, "+
			"COALESCE((SELECT COUNT(*) FROM likes JOIN posts ON likes.post_id = posts.id WHERE posts.user_id = users.id), 0) as counts_likes, "+
			"EXISTS(SELECT 1 FROM follows WHERE user_id = ? AND followed_id = users.id) as you_follow, "+
			"EXISTS(SELECT 1 FROM follows WHERE user_id = users.id AND followed_id = ?) as follows_you",
			curUserID, curUserID).
		Where("users.username = ?", username).
		First(&dbUser).Error; err != nil {
		return err
	}

	userResp := UserResponse{
		User: dbUser.User,
		Counts: Counts{
			Followers: dbUser.CountsFollowers,
			Following: dbUser.CountsFollowing,
			Posts:     dbUser.CountsPosts,
			Likes:     dbUser.CountsLikes,
		},
	}

	userResp.YouFollow = dbUser.YouFollow
	userResp.FollowsYou = dbUser.FollowsYou

	// Respond with the user data
	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    userResp,
	})
}
