package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func ListPosts(c *fiber.Ctx) error {
	var session *models.Session
	if c.Locals("session") != nil {
		session = c.Locals("session").(*models.Session)
	}

	var posts = make([]models.Post, 0)

	if err := lib.DB.
		Preload("User").
		Order("created_at desc").
		Where("type = ? OR type = ?", models.PostTypePost, models.PostTypeRepost).
		Find(&posts).Error; err != nil {
		return err
	}

	if session != nil {
		for i, post := range posts {
			// check if liked
			var like models.Like
			if err := lib.DB.
				Where(&models.Like{
					UserID: session.Connection.User.ID,
					PostID: post.ID,
				}).First(&like).Error; err != nil {
				posts[i].Liked = false
			} else {
				posts[i].Liked = true
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    posts,
	})
}

func GetPost(c *fiber.Ctx) error {
	var post models.Post
	if err := lib.DB.
		Preload("User").
		Where("id = ?", c.Params("post")).
		First(&post).Error; err != nil {
		return err
	}

	// isolate so it's only replies to this post
	var replies = make([]models.Post, 0)
	if err := lib.DB.
		Preload("User").
		Where("parent_id = ? and type = ?", post.ID, models.PostTypeReply).
		Find(&replies).Error; err != nil {
		return err
	}

	// replace all posts to only replies
	post.Posts = replies

	// check if liked
	var like models.Like
	if err := lib.DB.
		Where(&models.Like{
			UserID: c.Locals("session").(models.Session).Connection.User.ID,
			PostID: post.ID,
		}).First(&like).Error; err != nil {
		post.Liked = false
	} else {
		post.Liked = true
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    post,
	})
}

func GetPostsByUser(c *fiber.Ctx) error {
	var session = c.Locals("session").(models.Session)

	var user models.User
	if err := lib.DB.Where(&models.User{
		Username: c.Params("user"),
	}).First(&user).First(&user).Error; err != nil {
		return err
	}

	var posts = make([]models.Post, 0)
	if err := lib.DB.
		Where(&models.Post{
			UserID: user.ID,
		}).
		Preload("User").
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		return err
	}

	// check if liked
	for i, post := range posts {
		var like models.Like
		if err := lib.DB.
			Where(&models.Like{
				UserID: session.Connection.User.ID,
				PostID: post.ID,
			}).First(&like).Error; err != nil {
			posts[i].Liked = false
		} else {
			posts[i].Liked = true
		}
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    posts,
	})
}
