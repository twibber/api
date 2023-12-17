package account

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/img"
	"github.com/twibber/api/lib"
	"golang.org/x/exp/slices"
	"mime"
	"path/filepath"
	"strings"
)

func UpdateProfileImages(c *fiber.Ctx) error {
	session := lib.GetSession(c)

	uploadType := c.Params("type")

	allowedTypes := []string{"avatar", "banner"}
	if !slices.Contains(allowedTypes, uploadType) {
		return lib.NewError(fiber.StatusBadRequest, "Invalid upload type", nil)
	}

	file, err := c.FormFile(uploadType)
	if err != nil {
		return lib.NewError(fiber.StatusBadRequest, "Invalid file", nil)
	}

	// Check if file size exceeds 100MB
	const maxFileSize = 100 << 20 // 100MB in bytes
	if file.Size > maxFileSize {
		return lib.NewError(fiber.StatusBadRequest, "File size exceeds 100MB", nil)
	}

	// Check if file is an image
	fileExt := filepath.Ext(file.Filename)
	mimeType := mime.TypeByExtension(fileExt)
	if !strings.HasPrefix(mimeType, "image/") {
		return lib.NewError(fiber.StatusBadRequest, "File is not an image", nil)
	}

	// Upload the file to R2 and get the URL
	url, err := img.UploadFile(file, uploadType, session.Connection.User.ID)
	if err != nil {
		return lib.NewError(fiber.StatusInternalServerError, "Failed to upload file", &lib.ErrorDetails{
			Debug: err,
		})
	}

	if uploadType == "avatar" {
		session.Connection.User.Avatar = url
	} else {
		session.Connection.User.Banner = url
	}

	if err := lib.DB.Save(&session.Connection.User).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lib.BlankSuccess)
}
