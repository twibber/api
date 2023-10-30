package models

import "time"

// Session represents an authenticated session related to a connection.
type Session struct {
	ID string `gorm:"primaryKey" json:"id"`

	ConnectionID string      `gorm:"not null" json:"-"`
	Connection   *Connection `gorm:"foreignKey:ConnectionID;references:ID;constraint:OnDelete:CASCADE" json:"connection,omitempty"`

	ExpiresAt time.Time `gorm:"not null" json:"expires_at"` // Expiry date-time of the session

	// Timestamps
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}
