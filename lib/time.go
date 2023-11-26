package lib

import (
	"github.com/twibber/api/models"
	"time"
)

func NewDBTime() models.Timestamps {
	var now = time.Now()

	return models.Timestamps{
		CreatedAt: now,
		UpdatedAt: now,
	}
}
