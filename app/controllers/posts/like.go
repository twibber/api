package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

func LikePost(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var post models.Post
	if err := lib.DB.Where(&models.Post{
		ID: c.Params("post"),
	}).First(&post).Error; err != nil {
		return err
	}

	// check if like exists
	var like models.Like
	if err := lib.DB.Where(&models.Like{
		UserID: session.Connection.User.ID,
		PostID: post.ID,
	}).First(&like).Error; err == nil {
		return lib.NewError(fiber.StatusBadRequest, "You have already liked this post", nil)
	}

	if err := lib.DB.Create(&models.Like{
		ID:         utils.UUIDv4(),
		UserID:     session.Connection.User.ID,
		PostID:     post.ID,
		Timestamps: lib.NewDBTime(),
	}).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(lib.Response{
		Success: true,
		Data:    post,
	})
}

func UnlikePost(c *fiber.Ctx) error {
	session := c.Locals("session").(models.Session)

	var post models.Post
	if err := lib.DB.Where(&models.Post{
		ID: c.Params("post"),
	}).First(&post).Error; err != nil {
		return err
	}

	// check if like exists
	var like models.Like
	if err := lib.DB.Where(&models.Like{
		UserID: session.Connection.User.ID,
		PostID: post.ID,
	}).First(&like).Error; err != nil {
		return lib.NewError(fiber.StatusBadRequest, "You have not liked this post", nil)
	}

	return c.Status(fiber.StatusNoContent).JSON(lib.Response{
		Success: true,
		Data:    post,
	})
}
