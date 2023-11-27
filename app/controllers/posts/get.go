package posts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/models"
)

type PostQueryResult struct {
	models.Post
	models.User
	CountsLikes   int64 `gorm:"column:counts_likes"`
	CountsReplies int64 `gorm:"column:counts_replies"`
	CountsReposts int64 `gorm:"column:counts_reposts"`
	Liked         bool  `gorm:"column:liked"`
}

type PostResponse struct {
	models.Post
	Counts Counts `json:"counts"`
	Liked  bool   `json:"liked"`
}

type Counts struct {
	Likes   int64 `json:"likes"`
	Replies int64 `json:"replies"`
	Reposts int64 `json:"reposts"`
}

// ListPosts returns a list of all posts on the platform.
func ListPosts(c *fiber.Ctx) error {
	session := lib.GetSession(c)

	var userID string
	if session != nil {
		userID = session.Connection.User.ID
	}

	var posts []models.Post
	err := lib.DB.
		Model(&models.Post{}).
		Preload("User").
		Preload("Likes").
		Where("type IN ?", []string{string(models.PostTypePost), string(models.PostTypeRepost)}).
		Order("created_at DESC").
		Find(&posts).Error
	if err != nil {
		return err
	}

	var postResponses = make([]PostResponse, 0)
	for _, post := range posts {
		postResponses = append(postResponses, populatePostResponse(post, userID))
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    postResponses,
	})
}

// GetPostsByUser returns a list of posts by a user, in the same format as ListPosts.
func GetPostsByUser(c *fiber.Ctx) error {
	username := c.Params("user")

	session := lib.GetSession(c)

	var userID string
	if session != nil {
		userID = session.Connection.User.ID
	}

	var posts []models.Post
	err := lib.DB.
		Preload("User").
		Preload("Likes").
		Joins("JOIN users ON users.id = posts.user_id").
		Where("posts.type IN ? AND users.username = ?", []string{string(models.PostTypePost), string(models.PostTypeRepost)}, username).
		Order("posts.created_at DESC").
		Find(&posts).Error
	if err != nil {
		return err
	}

	var postResponses = make([]PostResponse, 0)
	for _, post := range posts {
		postResponses = append(postResponses, populatePostResponse(post, userID))
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    postResponses, // Ensure you return the postResponses, not the posts
	})
}

// GetPost returns a single post by ID with the Post.Posts attribute being filled with nothing but replies of the post sorted by order of posted, newer are further up.
func GetPost(c *fiber.Ctx) error {
	session := lib.GetSession(c)

	postID := c.Params("post")

	var userID string
	if session != nil {
		userID = session.Connection.User.ID
	}

	var post models.Post
	err := lib.DB.
		Model(&models.Post{}).
		Preload("User").
		Preload("Likes").
		Preload("Posts", "type = ?", models.PostTypeReply).
		Where("id = ?", postID).
		First(&post).Error
	if err != nil {
		return lib.ErrNotFound
	}

	// Populate the counts and liked status for the post
	postResponse := populatePostResponse(post, userID)

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    postResponse,
	})
}

func populatePostResponse(post models.Post, userID string) PostResponse {
	// Count likes, replies, and reposts
	likesCount := len(post.Likes)
	repliesCount := 0
	repostsCount := 0
	for _, p := range post.Posts {
		switch p.Type {
		case models.PostTypeReply:
			repliesCount++
		case models.PostTypeRepost:
			repostsCount++
		}
	}

	// Check if the post was liked by the current user
	likedByUser := false
	for _, like := range post.Likes {
		if like.UserID == userID {
			likedByUser = true
			break
		}
	}

	return PostResponse{
		Post: post,
		Counts: Counts{
			Likes:   int64(likesCount),
			Replies: int64(repliesCount),
			Reposts: int64(repostsCount),
		},
		Liked: likedByUser,
	}
}
