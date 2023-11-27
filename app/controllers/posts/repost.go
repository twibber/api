package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type RepostDTO struct {
	Content string `json:"content" validate:"omitempty,max=512"`
}

func CreateRepost(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	postID := c.Params("post")

	var dto RepostDTO
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	var parentPost models.Post
	if err := lib.DB.Where("id = ?", postID).First(&parentPost).Error; err != nil {
		return err
	}

	if parentPost.Type == models.PostTypeRepost {
		return lib.NewError(fiber.StatusBadRequest, "You cannot repost a repost", nil)
	}

	dbReply := &models.Post{
		ID:         utils.UUIDv4(),
		UserID:     session.Connection.User.ID,
		Type:       models.PostTypeRepost,
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
