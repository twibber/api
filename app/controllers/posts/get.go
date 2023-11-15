package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func ListPosts(c *fiber.Ctx) error {
	var posts = make([]models.Post, 0)

	// get all posts sorted by date posted
	if err := lib.DB.
		Preload("User").
		Preload("Posts").
		Preload("Posts.User").
		Order("created_at desc").
		Find(&posts).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    posts,
	})
}
