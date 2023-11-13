package models

import (
	"time"
)

// User represents the system user with related authentication details and relationships.
type User struct {
	ID string `gorm:"primaryKey" json:"id"`

	// Details
	Username    string `gorm:"size:255;not null" json:"username"`
	DisplayName string `gorm:"size:255" json:"display_name"`

	Admin bool `gorm:"not null;default:false" json:"admin,omitempty"`

	// Auth
	Email string `gorm:"size:255;unique;not null" json:"email"`
	MFA   string `json:"-"` // Hidden: Multi-Factor Authentication

	// Account Suspension
	Suspended bool `gorm:"default:false" json:"suspended,omitempty"`

	// Relationships
	Connections []Connection `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"connections,omitempty"`
	Posts       []Post       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"posts,omitempty"`
	Likes       []Like       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"likes,omitempty"`

	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}
