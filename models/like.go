package models

import "time"

type Like struct {
	ID string `json:"id" gorm:"primaryKey"`

	// User
	UserID string `json:"-"`
	User   *User  `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`

	// Post
	PostID string `json:"-"`
	Post   *Post  `json:"post,omitempty" gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`

	// Timestamps
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoCreateTime"`
}
