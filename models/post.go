package models

import "time"

type Post struct {
	ID string `json:"id" gorm:"primaryKey"`

	AuthorID string `json:"author_id" gorm:"index"`

	// Timestamps
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoCreateTime"`
}
