package models

import "time"

type Follow struct {
	ID string `json:"id" gorm:"primaryKey"`

	// User
	UserID string `json:"user_id" gorm:"index"`

	// Follower
	FollowerID string `json:"-" gorm:"index"`
	Follower   *User  `json:"follower,omitempty" gorm:"foreignKey:FollowerID;references:ID;constraint:OnDelete:CASCADE"`

	// Timestamps
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoCreateTime"`
}
