package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
	"time"
)

type CreatePostDTO struct {
	Content string `json:"content" validate:"required"`
}

func CreatePost(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var post CreatePostDTO
	if err := lib.ParseAndValidate(c, &post); err != nil {
		return err
	}

	if err := lib.DB.Create(&models.Post{
		ID:        utils.UUIDv4(),
		UserID:    session.Connection.User.ID,
		Type:      models.PostTypePost,
		Content:   &post.Content,
		Likes:     []models.Like{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(lib.Response{
		Success: true,
		Data:    post,
	})
}
