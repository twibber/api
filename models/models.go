package models

import "time"

var Models = []any{
	&User{},
	&Connection{},
	&Session{},
	&Post{},
	&Like{},
	&Follow{},
}

// Timestamps is an embedded struct that contains the created and updated times for a model.
type Timestamps struct {
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`
}
