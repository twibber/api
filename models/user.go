package models

// User represents the system user with related authentication details and profile information.
type User struct {
	ID     string `gorm:"primaryKey" json:"id"`                     // Unique identifier for the user
	JoinID int64  `gorm:"not null;default:0;unique" json:"join_id"` // A unique joining ID for the user

	Username    string `gorm:"size:255;not null;unique" json:"username"` // The user's chosen username, unique across the system
	DisplayName string `gorm:"size:255" json:"display_name"`             // The user's display name, shown to other users

	Avatar string `gorm:"size:255" json:"avatar"` // URL to the user's avatar image
	Banner string `gorm:"size:255" json:"banner"` // URL to the user's banner image

	Admin          bool `gorm:"not null;default:false" json:"admin"`           // Flag indicating whether the user has administrative privileges
	VerifiedPerson bool `gorm:"not null;default:false" json:"verified_person"` // Flag indicating whether the user is a verified person

	Email string `gorm:"size:255;unique;not null" json:"-"` // The user's email address, hidden in JSON responses

	MFA       string `json:"-"`                              // Multi-Factor Authentication details, if enabled, not exposed through API
	Suspended bool   `gorm:"default:false" json:"suspended"` // Flag indicating whether the user's account is suspended

	// Relationships
	Following []Follow `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"following,omitempty"`     // List of users that this user is following
	Followers []Follow `gorm:"foreignKey:FollowedID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"followers,omitempty"` // List of users that follow this user

	Connections []Connection `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"connections,omitempty"` // Authentication connections associated with the user

	Posts []Post `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Posts created by the user
	Likes []Like `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"` // Likes made by the user on posts

	Timestamps // Embedded struct for created and updated timestamps
}

// Follow represents a relationship where a User is following another User.
type Follow struct {
	ID string `gorm:"primaryKey" json:"id"` // Unique identifier for the follow relationship

	UserID string `gorm:"not null" json:"user_id"`                                                                             // ID of the user who is following
	User   User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"` // The user who is following

	FollowedID string `gorm:"not null" json:"followed_id"`                                                                                 // ID of the user being followed
	Followed   User   `gorm:"foreignKey:FollowedID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"followed,omitempty"` // The user being followed

	Timestamps // Embedded struct for created and updated timestamps
}
