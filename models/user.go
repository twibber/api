package models

import "time"

type User struct {
	ID string `json:"id" gorm:"primaryKey"`

	// Details
	Username string `json:"username"`

	// Auth
	Email    string `json:"email"`
	Verified bool   `json:"verified,omitempty"`
	MFA      string `json:"-"`

	// Account Suspension
	// this limits the account's actions, this is in place in case a suspicious account needs to be kept on record.
	Suspended bool `json:"suspended,omitempty"`

	// Authentication
	Connections []Connection `json:"connections,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Sessions    []Session    `json:"sessions,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Timestamps
	CreatedAt time.Time `json:"-" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"-" gorm:"autoCreateTime"`
}
