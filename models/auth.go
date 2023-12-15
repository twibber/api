package models

import (
	"github.com/lib/pq"
	"time"
)

// ConnectionType represents the type of authentication method.
type ConnectionType string

// Predefined constants for ConnectionType.
const (
	EmailType  ConnectionType = "email"
	GoogleType ConnectionType = "google"
	GitHubType ConnectionType = "github"
)

func (c ConnectionType) WithID(id string) string {
	return string(c) + ":" + id
}

// Connection represents the authentication connections related to a user.
type Connection struct {
	BaseModel

	UserID     string    `gorm:"not null" json:"-"`
	User       *User     `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	TOTPVerify string    `json:"-"` // TOTP Verification code, not exposed through API
	Password   string    `json:"-"` // Password for the connection, not exposed through API
	Verified   bool      `gorm:"default:false" json:"verified"`
	Sessions   []Session `gorm:"foreignKey:ConnectionID;references:ID;constraint:OnDelete:CASCADE" json:"sessions,omitempty"`
}

// Session represents an authenticated session related to a connection.
type Session struct {
	BaseModel

	ConnectionID string      `gorm:"not null" json:"-"`
	Connection   *Connection `gorm:"foreignKey:ConnectionID;references:ID;constraint:OnDelete:CASCADE" json:"connection,omitempty"`
	Info         SessionInfo `gorm:"embedded;embeddedPrefix:info_" json:"info,omitempty"`
	ExpiresAt    time.Time   `gorm:"not null" json:"expires_at"`
}

// SessionInfo holds information about the session such as IP address and user agent.
type SessionInfo struct {
	IPAddresses pq.StringArray `gorm:"type:text[]" json:"ip_addresses,omitempty"`
	UserAgent   string         `gorm:"size:255" json:"user_agent,omitempty"`
}
