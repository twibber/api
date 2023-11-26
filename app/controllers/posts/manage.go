package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
	"time"
)

type CreatePostDTO struct {
	Content string `json:"content" validate:"required,max=512,min=1"`
}

func CreatePost(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var post CreatePostDTO
	if err := lib.ParseAndValidate(c, &post); err != nil {
		return err
	}

	if err := lib.DB.Create(&models.Post{
		ID:         utils.UUIDv4(),
		UserID:     session.Connection.User.ID,
		Type:       models.PostTypePost,
		Content:    &post.Content,
		Timestamps: lib.NewDBTime(),
	}).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(lib.Response{
		Success: true,
		Data:    post,
	})
}

func DeletePost(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var selector = &models.Post{
		ID:     c.Params("post"),
		UserID: session.Connection.User.ID,
	}

	if session.Connection.User.Admin {
		selector.UserID = ""
	}

	var post models.Post
	if err := lib.DB.Where(selector).First(&post).Error; err != nil {
		return err
	}

	if !session.Connection.User.Admin && time.Since(post.CreatedAt) > time.Minute*5 {
		return lib.NewError(fiber.StatusBadRequest, "You cannot delete a post after more than 5 minutes.", nil)
	}

	if err := lib.DB.Delete(&post).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
