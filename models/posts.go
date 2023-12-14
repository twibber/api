package models

// PostType represents the type of the post.
type PostType string

// Constants for different post types.
const (
	PostTypePost   PostType = "post"
	PostTypeReply  PostType = "reply"
	PostTypeRepost PostType = "repost"
)

// Post represents a user's post with potential relationships to other posts.
type Post struct {
	BaseModel

	UserID string `gorm:"not null" json:"user_id"`                                                                             // ID of the user who created the post
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"` // The user who created the post

	ParentID *string `gorm:"index" json:"parent_id,omitempty"`                                                                         // ID of the parent post, if this is a reply or repost
	Parent   *Post   `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent,omitempty"` // The parent post

	Type    PostType `json:"type"` // The type of the post (post, reply, repost)
	Content *string  `gorm:"type:text" json:"content,omitempty"`

	Posts []Post `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"posts,omitempty"` // Posts associated with the post
	Likes []Like `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"likes,omitempty"`   // Likes associated with the post

	// Ignored by GORM and populated by the handler.
	Liked bool `gorm:"-" json:"liked,omitempty"` // Flag indicating whether the post was liked by the current user

	// Counts are ignored by GORM and are populated by the handler.
	Counts struct {
		Likes   int `gorm:"-" json:"likes"`   // Number of likes on the post
		Replies int `gorm:"-" json:"replies"` // Number of replies to the post
		Reposts int `gorm:"-" json:"reposts"` // Number of reposts of the post
	} `gorm:"-" json:"counts,omitempty"` // Counts associated with the post
}

// Like represents a 'like' given by a user to a post.
type Like struct {
	BaseModel

	UserID string `gorm:"not null" json:"user_id"`                                                                             // ID of the user who liked the post
	User   *User  `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"` // The user who liked the post

	PostID string `gorm:"not null" json:"post_id"`                                                                             // ID of the post that was liked
	Post   *Post  `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post,omitempty"` // The post that was liked
}
