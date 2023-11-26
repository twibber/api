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

	var queryResults []PostQueryResult
	query := `
SELECT p.*, 
    u.*,
    COALESCE(l.likes_count, 0) as "counts_likes",
    COALESCE(r.replies_count, 0) as "counts_replies",
    COALESCE(rp.reposts_count, 0) as "counts_reposts",
    EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = p.id) as "liked"
FROM posts p
JOIN users u ON p.user_id = u.id
LEFT JOIN (SELECT post_id, COUNT(*) as likes_count FROM likes GROUP BY post_id) l ON l.post_id = p.id
LEFT JOIN (SELECT parent_id, COUNT(*) as replies_count FROM posts WHERE type = 'reply' GROUP BY parent_id) r ON r.parent_id = p.id
LEFT JOIN (SELECT parent_id, COUNT(*) as reposts_count FROM posts WHERE type = 'repost' GROUP BY parent_id) rp ON rp.parent_id = p.id
WHERE p.type IN ('post', 'repost')
ORDER BY p.created_at DESC
`
	// Pass the sessionUserID as an argument to the raw SQL query
	if err := lib.DB.Raw(query, userID).Scan(&queryResults).Error; err != nil {
		return err
	}

	// Convert the QueryResult to PostResponse
	var respPosts []PostResponse
	for _, qr := range queryResults {
		curPost := PostResponse{
			Post: qr.Post,
			Counts: Counts{
				Likes:   qr.CountsLikes,
				Replies: qr.CountsReplies,
				Reposts: qr.CountsReposts,
			},
			Liked: qr.Liked,
		}

		// Assign the user data to the Post
		curPost.Post.User = qr.User

		respPosts = append(respPosts, curPost)
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    respPosts,
	})
}

// GetPostsByUser returns a list of posts by a user, in the same format as ListPosts.
func GetPostsByUser(c *fiber.Ctx) error {
	username := c.Params("user")
	session := lib.GetSession(c)
	userID := ""

	if session != nil {
		userID = session.Connection.User.ID
	}

	var queryResults []PostQueryResult
	query := `SELECT p.*, 
		u.*,
		COALESCE(l.likes_count, 0) as "counts_likes",
		COALESCE(r.replies_count, 0) as "counts_replies",
		COALESCE(rp.reposts_count, 0) as "counts_reposts",
		EXISTS(SELECT 1 FROM likes WHERE user_id = ? AND post_id = p.id) as "liked"
	FROM posts p
	JOIN users u ON p.user_id = u.id
	LEFT JOIN (SELECT post_id, COUNT(*) as likes_count FROM likes GROUP BY post_id) l ON l.post_id = p.id
	LEFT JOIN (SELECT parent_id, COUNT(*) as replies_count FROM posts WHERE type = 'reply' GROUP BY parent_id) r ON r.parent_id = p.id
	LEFT JOIN (SELECT parent_id, COUNT(*) as reposts_count FROM posts WHERE type = 'repost' GROUP BY parent_id) rp ON rp.parent_id = p.id
	WHERE u.username = ? AND p.type IN ('post', 'repost')
	ORDER BY p.created_at DESC`

	// Pass the sessionUserID as an argument to the raw SQL query
	if err := lib.DB.Raw(query, userID, username).Scan(&queryResults).Error; err != nil {
		return err
	}

	// Convert the QueryResult to PostResponse
	var respPosts []PostResponse
	for _, qr := range queryResults {
		curPost := PostResponse{
			Post: qr.Post,
			Counts: Counts{
				Likes:   qr.CountsLikes,
				Replies: qr.CountsReplies,
				Reposts: qr.CountsReposts,
			},
			Liked: qr.Liked,
		}

		// Assign the user data to the Post
		curPost.Post.User = qr.User

		respPosts = append(respPosts, curPost)
	}

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    respPosts,
	})
}

// GetPost returns a single post by ID with the Post.Posts attribute being filled with nothing but replies of the post sorted by order of posted, newer are further up.
func GetPost(c *fiber.Ctx) error {
	postID := c.Params("post") // The parameter for PostID is :post
	var session *models.Session
	var sessionUserID string

	if s, ok := c.Locals("session").(*models.Session); ok {
		session = s
		sessionUserID = session.Connection.User.ID
	}

	// Get the post with user data, like count, and whether the current user has liked the post.
	var postResult PostQueryResult
	postQuery := `
        SELECT p.*,
            u.*,
            COALESCE(l.likes_count, 0) as "counts_likes",
            COALESCE(r.replies_count, 0) as "counts_replies",
            COALESCE(rp.reposts_count, 0) as "counts_reposts"
        FROM posts p
        JOIN users u ON p.user_id = u.id
        LEFT JOIN (SELECT post_id, COUNT(*) as likes_count FROM likes GROUP BY post_id) l ON l.post_id = p.id
        LEFT JOIN (SELECT parent_id, COUNT(*) as replies_count FROM posts WHERE type = 'reply' GROUP BY parent_id) r ON r.parent_id = p.id
        LEFT JOIN (SELECT parent_id, COUNT(*) as reposts_count FROM posts WHERE type = 'repost' GROUP BY parent_id) rp ON rp.parent_id = p.id
        WHERE p.id = ?
    `

	if err := lib.DB.Raw(postQuery, postID).Scan(&postResult).Error; err != nil {
		return err
	}

	// Now, get all replies for this post.
	var replies []models.Post
	repliesQuery := `
        SELECT * FROM posts
        WHERE parent_id = ? AND type = 'reply'
        ORDER BY created_at DESC
    `
	if err := lib.DB.Raw(repliesQuery, postID).Scan(&replies).Error; err != nil {
		return err
	}

	// Determine if the current session user has liked the post
	liked := false
	if sessionUserID != "" {
		var like models.Like
		err := lib.DB.
			Where("user_id = ? AND post_id = ?", sessionUserID, postID).
			First(&like).Error
		liked = err == nil
	}

	// Construct the response
	postResponse := PostResponse{
		Post: postResult.Post,
		Counts: Counts{
			Likes:   postResult.CountsLikes,
			Replies: postResult.CountsReplies,
			Reposts: postResult.CountsReposts,
		},
		Liked: liked,
	}
	postResponse.Post.User = postResult.User // Assign user data to the Post
	postResponse.Post.Posts = replies        // Assign the replies to the Posts field

	return c.Status(fiber.StatusOK).JSON(lib.Response{
		Success: true,
		Data:    postResponse,
	})
}
