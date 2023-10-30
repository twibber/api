package models

import "time"

type Type string

const (
	Email  Type = "email"
	Google Type = "google"
)

func (t Type) WithID(id string) string {
	return string(t) + ":" + id
}

// Connection represents the authentication connections related to a user.
type Connection struct {
	ID string `gorm:"primaryKey" json:"id"` // e.g., google:id or github:id

	UserID string `gorm:"not null" json:"-"`
	User   *User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`

	TOTPVerify string `json:"-"` // Hidden: TOTP Verification code
	Password   string `json:"-"` // Hidden: Password for the connection

	Verified bool `gorm:"default:false" json:"verified"`

	Sessions []Session `gorm:"foreignKey:ConnectionID;references:ID;constraint:OnDelete:CASCADE" json:"sessions,omitempty"`

	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}
