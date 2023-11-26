package lib

import (
	"github.com/twibber/api/models"
	"time"
)

func NewDBTime() models.Timestamps {
	return models.Timestamps{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
