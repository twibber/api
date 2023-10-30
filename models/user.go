package models

import (
	"time"
)

// Level represents the access level of a user in the system.
type Level string

// User levels
const (
	Admin   Level = "admin"
	Support Level = "support"
	Member  Level = "user"
)

// User represents the system user with related authentication details and relationships.
type User struct {
	ID string `gorm:"primaryKey" json:"id"`

	// Details
	Username string `gorm:"size:255;not null" json:"username"`

	Level Level `gorm:"not null;default:user" json:"level"`

	// Auth
	Email string `gorm:"size:255;unique;not null" json:"email"`
	MFA   string `json:"-"` // Hidden: Multi-Factor Authentication

	// Account Suspension
	Suspended bool `gorm:"default:false" json:"suspended,omitempty"`

	// Relationships
	Connections []Connection `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"connections,omitempty"`

	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}
