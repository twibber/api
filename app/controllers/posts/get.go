package posts

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
	"net/url"
	"regexp"
	"strings"
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
	replaceImageURLsWithProxy(post.Content)

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
		replaceImageURLsWithProxy(subPost.Content)

		switch subPost.Type {
		case models.PostTypeReply:
			post.Counts.Replies++
			if includeReplies {
				populatePostCounts(&subPost, userID, false)
				replaceImageURLsWithProxy(subPost.Content)
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
		replaceImageURLsWithProxy(post.Parent.Content)
		populatePostCounts(post.Parent, userID, false)
	}
}

func replaceImageURLsWithProxy(content *string) {
	keyBin, err := hex.DecodeString(lib.Config.ImgproxyKey)
	if err != nil {
		log.Fatal("Error decoding key: ", err)
	}

	saltBin, err := hex.DecodeString(lib.Config.ImgproxySalt)
	if err != nil {
		log.Fatal("Error decoding salt: ", err)
	}

	re := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	sections := strings.Split(*content, "```")

	for i, section := range sections {
		if i%2 == 0 {
			section = re.ReplaceAllStringFunc(section, func(url string) string {
				matches := re.FindStringSubmatch(url)
				if len(matches) < 2 {
					return url
				}

				originalURL := matches[1]
				signedURL := signURL(originalURL, keyBin, saltBin)
				return strings.Replace(url, originalURL, signedURL, 1)
			})
			sections[i] = section
		}
	}

	*content = strings.Join(sections, "```")
}

func signURL(imgURL string, keyBin []byte, saltBin []byte) string {
	escapedImgURL := url.QueryEscape(imgURL)

	// Set the desired quality, between 0-100 (lower means more compression).
	quality := 80

	// Set the desired width and height based on your UI's requirements.
	width := 200
	height := 200

	// Choose an image format, such as WebP for better compression.
	format := "webp"

	// Device Pixel Ratio, use 1 for standard resolution, higher for high-DPI displays.
	dpr := 1

	path := fmt.Sprintf("/rs:fill:%d:%d:0/q:%d/dpr:%d/f:%s/plain/%s", width, height, quality, dpr, format, escapedImgURL)

	// Calculate the HMAC digest
	mac := hmac.New(sha256.New, keyBin)
	mac.Write(saltBin)      // Writing salt first
	mac.Write([]byte(path)) // Writing the path
	signature := mac.Sum(nil)

	// Base64 URL-Safe Encoding of the signature
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)
	encodedSignature = strings.TrimRight(encodedSignature, "=")

	// Construct the final signed URL
	return fmt.Sprintf("%s/%s%s", lib.Config.ImgproxyURL, encodedSignature, path)
}
