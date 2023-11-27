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
	userID := ""

	if session != nil {
		userID = session.Connection.User.ID
	}

	var respPosts []PostQueryResult
	if err := lib.DB.
		Table("posts").
		Select("posts.*, users.*, "+ // Select all columns from posts and users
			"COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = posts.id), 0) as counts_likes, "+ // Count likes
			"COALESCE((SELECT COUNT(*) FROM posts WHERE parent_id = posts.id AND type = 'reply'), 0) as counts_replies, "+ // Count replies
			"COALESCE((SELECT COUNT(*) FROM posts WHERE parent_id = posts.id AND type = 'repost'), 0) as counts_reposts, "+ // Count reposts
			"EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = posts.id) as liked", // Check if the current user has liked the post
									userID). // Pass the current user ID to the query
		Joins("JOIN users ON posts.user_id = users.id").      // Join the users table
		Where("posts.type IN ?", []string{"post", "repost"}). // Only select posts and reposts
		Order("posts.created_at DESC").                       // Order by newest first
		Scan(&respPosts).Error; err != nil {
		return err
	}

	var posts = make([]PostResponse, 0)
	for _, post := range respPosts {
		var curPost = PostResponse{
			Post: post.Post,
			Counts: Counts{
				Likes:   post.CountsLikes,
				Replies: post.CountsReplies,
				Reposts: post.CountsReposts,
			},
			Liked: post.Liked,
		}

		curPost.Post.User = post.User

		posts = append(posts, curPost)
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    posts,
	})
}

// GetPostsByUser returns a list of posts by a user, in the same format as ListPosts.
func GetPostsByUser(c *fiber.Ctx) error {
	session := lib.GetSession(c)
	userID := ""

	if session != nil {
		userID = session.Connection.User.ID
	}

	var respPosts []PostQueryResult
	err := lib.DB.
		Table("posts").
		Select("posts.*, users.*, "+ // Select all columns from posts and users
			"COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = posts.id), 0) as counts_likes, "+ // Count likes
			"COALESCE((SELECT COUNT(*) FROM posts WHERE parent_id = posts.id AND type = 'reply'), 0) as counts_replies, "+ // Count replies
			"COALESCE((SELECT COUNT(*) FROM posts WHERE parent_id = posts.id AND type = 'repost'), 0) as counts_reposts, "+ // Count reposts
			"EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = posts.id) as liked", // Check if the current user has liked the post
			userID).
		Joins("JOIN users ON posts.user_id = users.id").
		Where("posts.type IN ? AND users.username = ?", []string{"post", "repost"}, c.Params("user")).
		Order("posts.created_at DESC").
		Scan(&respPosts).Error
	if err != nil {
		return err
	}

	var posts = make([]PostResponse, 0)
	for _, post := range respPosts {
		var curPost = PostResponse{
			Post: post.Post,
			Counts: Counts{
				Likes:   post.CountsLikes,
				Replies: post.CountsReplies,
				Reposts: post.CountsReposts,
			},
			Liked: post.Liked,
		}

		curPost.Post.User = post.User

		posts = append(posts, curPost)
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    posts,
	})
}

// GetPost returns a single post by ID with the Post.Posts attribute being filled with nothing but replies of the post sorted by order of posted, newer are further up.
func GetPost(c *fiber.Ctx) error {
	session := lib.GetSession(c)
	userID := ""

	if session != nil {
		userID = session.Connection.User.ID
	}

	var respPost PostQueryResult
	err := lib.DB.
		Select("posts.*, users.*, "+ // Select all columns from posts and users
			"COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = posts.id), 0) as counts_likes, "+ // Count likes
			"COALESCE((SELECT COUNT(*) FROM posts WHERE parent_id = posts.id AND type = 'reply'), 0) as counts_replies, "+ // Count replies
			"COALESCE((SELECT COUNT(*) FROM posts WHERE parent_id = posts.id AND type = 'repost'), 0) as counts_reposts, "+ // Count reposts
			"EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = posts.id) as liked", // Check if the current user has liked the post
			userID).
		Joins("JOIN users ON posts.user_id = users.id").
		Where("posts.type IN ? AND users.username = ?", []string{"post", "repost"}, c.Params("user")).
		Order("posts.created_at DESC").
		First(&respPost).Error
	if err != nil {
		return err
	}

	var curPost = PostResponse{
		Post: respPost.Post,
		Counts: Counts{
			Likes:   respPost.CountsLikes,
			Replies: respPost.CountsReplies,
			Reposts: respPost.CountsReposts,
		},
		Liked: respPost.Liked,
	}

	curPost.Post.User = respPost.User

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    curPost,
	})
}
