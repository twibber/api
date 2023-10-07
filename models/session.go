package models

import "time"

type Session struct {
	ID string `json:"id" gorm:"primaryKey"`

	UserID string `json:"-"`
	User   *User  `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`

	ExpiresAt time.Time `json:"expires_at"`

	// Timestamps
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoCreateTime"`
}
