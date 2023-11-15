package models

import "time"

type PostType string

const (
	PostTypePost   PostType = "post"
	PostTypeReply  PostType = "reply"
	PostTypeRepost PostType = "repost"
)

type Post struct {
	ID string `gorm:"primaryKey" json:"id"`

	UserID string `gorm:"not null" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	ParentID *string `json:"parent_id"`
	Parent   *Post   `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"parent,omitempty"`

	Type    PostType `json:"type"`
	Content *string  `gorm:"type:text" json:"content"`

	Posts []Post `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"posts,omitempty"`
	Likes []Like `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"likes,omitempty"`

	// non database fields for replies and reposts to be used in ListPosts
	Replies []Post `gorm:"-" json:"replies,omitempty"`
	Reposts []Post `gorm:"-" json:"reposts,omitempty"`

	// Hidden database fields for count values for likes reposts and replies
	LikeCount   int `gorm:"-" json:"like_count"`
	RepostCount int `gorm:"-" json:"repost_count"`
	ReplyCount  int `gorm:"-" json:"reply_count"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}

type Like struct {
	ID string `gorm:"primaryKey" json:"id"`

	// Relationships
	UserID string `gorm:"not null" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	PostID string `gorm:"not null" json:"post_id"`
	Post   Post   `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"post,omitempty"`

	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}
