package lib

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/twibber/api/models"
	"gorm.io/gorm"
	"strings"
	"time"
)

func GetSession(c *fiber.Ctx) *models.Session {
	// Extracts auth token from header or cookie
	authHeader := c.Get("Authorization")
	authCookie := c.Cookies("Authorization")

	// Rejects request if no auth details are provided
	if authHeader == "" && authCookie == "" {
		logrus.Debug("No auth header or cookie")
		return nil
	}

	var authToken string

	// Parses auth token from header or uses cookie
	if authHeader != "" {
		authParsed := strings.Split(authHeader, " ")
		if len(authParsed) < 2 {
			logrus.Debug("Invalid auth header")
			return nil
		}
		authToken = authParsed[1]
	} else if authCookie != "" {
		authToken = authCookie
	}

	// Fetches session from database using auth token
	var session models.Session
	if err := DB.Where(models.Session{ID: authToken}).Preload("Connection").Preload("Connection.User").First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Debug("Session not found")
			return nil
		}
		return nil
	}

	if session.ExpiresAt.Before(time.Now()) {
		DB.Delete(&session)
		return nil
	}

	return &session
}
