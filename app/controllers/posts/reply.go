package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type ReplyDTO struct {
	Content string `json:"content" validate:"required,max=512,min=1,notblank"`
}

func CreateReply(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var dto ReplyDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	var parentPost models.Post
	if err := lib.DB.Where(&models.Post{
		BaseModel: models.BaseModel{ID: c.Params("post")},
	}).First(&parentPost).Error; err != nil {
		return err
	}

	dbReply := &models.Post{
		UserID:   session.Connection.User.ID,
		Type:     models.PostTypeReply,
		ParentID: &parentPost.ID,
		Content:  &dto.Content,
	}

	if err := lib.DB.Create(&dbReply).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(lib.Response{
		Success: true,
		Data:    dbReply,
	})
}
