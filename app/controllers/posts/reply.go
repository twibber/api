package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type ReplyDTO struct {
	Content string `json:"content" validate:"required,max=512"`
}

func CreateReply(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var dto ReplyDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	var parentPost models.Post
	if err := lib.DB.Where(&models.Post{
		ID: c.Params("post"),
	}).First(&parentPost).Error; err != nil {
		return err
	}

	dbReply := &models.Post{
		ID:         utils.UUIDv4(),
		UserID:     session.Connection.User.ID,
		Type:       models.PostTypeReply,
		ParentID:   &parentPost.ID,
		Content:    &dto.Content,
		Timestamps: lib.NewDBTime(),
	}

	if err := lib.DB.Create(&dbReply).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(lib.Response{
		Success: true,
		Data:    dbReply,
	})
}
