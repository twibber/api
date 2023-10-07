package models

import "time"

type Type string

const (
	Email  Type = "email"
	Google Type = "google"
	GitHub Type = "github"
)

type Connection struct {
	ID string `json:"id" gorm:"primaryKey"`

	UserID string `json:"-"`
	User   *User  `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`

	Type     Type   `json:"type"`
	Password string `json:"-"`

	// Timestamps
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoCreateTime"`
}
