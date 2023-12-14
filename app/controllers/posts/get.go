package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

// ListPosts returns a list of all posts on the platform.
func ListPosts(c *fiber.Ctx) error {
	session := lib.GetSession(c)
	userID := ""

	if session != nil {
		userID = session.Connection.User.ID
	}

	var posts []models.Post
	if err := lib.DB.
		Model(&models.Post{}).
		Preload("User").
		Preload("Likes").
		Preload("Posts").
		Preload("Parent").
		Preload("Parent.User").
		Preload("Parent.Likes").
		Preload("Parent.Posts").
		Preload("Parent.Parent").
		Where("type = ? OR type = ?", models.PostTypePost, models.PostTypeRepost).
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		return err
	}

	for i := range posts {
		populatePostCounts(&posts[i], userID, false)
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    posts,
	})
}

func GetPostsByUser(c *fiber.Ctx) error {
	username := c.Params("user")

	var user models.User
	if err := lib.DB.
		Model(&models.User{}).
		Where(&models.User{Username: username}).
		First(&user).Error; err != nil {
		return err
	}

	var posts []models.Post
	if err := lib.DB.
		Model(&models.Post{}).
		Preload("User").
		Preload("Likes").
		Preload("Posts").
		Preload("Parent").
		Preload("Parent.User").
		Preload("Parent.Likes").
		Preload("Parent.Posts").
		Preload("Parent.Parent").
		Where(&models.Post{
			UserID: user.ID,
		}).
		Where("type = ? OR type = ?", models.PostTypePost, models.PostTypeRepost).
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		return err
	}

	sessionUserID := ""
	if session := lib.GetSession(c); session != nil {
		sessionUserID = session.Connection.User.ID
	}

	for i := range posts {
		populatePostCounts(&posts[i], sessionUserID, false)
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    posts,
	})
}

func GetPost(c *fiber.Ctx) error {
	postID := c.Params("post")

	var post models.Post
	if err := lib.DB.
		Model(&models.Post{}).
		Preload("User").
		Preload("Likes").
		Preload("Posts").
		Preload("Parent").
		Preload("Parent.User").
		Preload("Parent.Likes").
		Preload("Parent.Posts").
		Preload("Parent.Parent").
		Preload("Posts.User").
		Preload("Posts.Likes").
		Preload("Posts.Posts").
		Where("id = ?", postID).
		First(&post).Error; err != nil {
		return err
	}

	sessionUserID := ""
	if session := lib.GetSession(c); session != nil {
		sessionUserID = session.Connection.User.ID
	}

	populatePostCounts(&post, sessionUserID, true)

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    post,
	})
}

// populatePostCounts populates the counts and liked fields on a post.
func populatePostCounts(post *models.Post, userID string, includeReplies bool) {
	// check if liked by user
	for _, like := range post.Likes {
		if like.UserID == userID {
			post.Liked = true
			break
		}
	}

	// count likes on post
	post.Counts.Likes = len(post.Likes)

	var replies []models.Post
	// count replies and reposts on post
	for _, subPost := range post.Posts {
		switch subPost.Type {
		case models.PostTypeReply:
			post.Counts.Replies++
			if includeReplies {
				populatePostCounts(&subPost, userID, false)
				replies = append(replies, subPost)
			}
		case models.PostTypeRepost:
			post.Counts.Reposts++
		}
	}

	if includeReplies {
		post.Posts = replies
	} else {
		post.Posts = nil
	}

	// if there is a parent post, recursively populate its counts
	if post.Parent != nil {
		populatePostCounts(post.Parent, userID, false)
	}
}
