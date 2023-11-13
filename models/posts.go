package models

import "time"

type PostType string

const (
	PostTypePost  PostType = "post"
	PostTypeQuote PostType = "quote"
	PostTypeReply PostType = "reply"
)

type Post struct {
	ID     string `gorm:"primaryKey" json:"id"`
	UserID string `gorm:"not null" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`

	ParentID *string `json:"parent_id"`
	Parent   *Post   `gorm:"foreignKey:ParentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"parent,omitempty"`

	Type    PostType `json:"type"`
	Content string   `gorm:"type:text;not null" json:"content"`

	Likes []Like `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"likes,omitempty"`

	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
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
