package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"
	"github.com/twibber/api/models"
	"gorm.io/gorm"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/twibber/api/lib"
	"golang.org/x/oauth2"
)

var (
	oauth2Configs = map[string]*oauth2.Config{
		"google": {
			ClientID:     lib.Config.GoogleClient,
			ClientSecret: lib.Config.GoogleSecret,
			RedirectURL:  fmt.Sprintf("%s/auth/oauth/google/callback", lib.Config.APIURL),
			Scopes:       []string{"profile", "email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
		},
	}
)

func AuthorisationURL(c *fiber.Ctx) error {
	provider := c.Params("provider")
	conf, ok := oauth2Configs[provider]
	if !ok {
		return lib.NewError(fiber.StatusBadRequest, "Invalid OAuth provider.", nil)
	}

	state := utils.UUIDv4()
	url := conf.AuthCodeURL(state)

	c.ClearCookie(fmt.Sprintf("oauth_%s_state", provider))
	c.Cookie(&fiber.Cookie{
		Name:     fmt.Sprintf("oauth_%s_state", provider),
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
	})

	return c.Redirect(url)
}

func OAuthCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	conf, ok := oauth2Configs[provider]
	if !ok {
		return lib.NewError(fiber.StatusBadRequest, "Invalid OAuth provider.", nil)
	}

	state := c.Query("state")
	cookieState := c.Cookies(fmt.Sprintf("oauth_%s_state", provider))

	log.Info(fmt.Sprintf("state: %s, cookieState: %s", state, cookieState))

	// Validate state token.
	if state == "" || state != cookieState {
		return lib.NewError(fiber.StatusBadRequest, "Invalid state token.", nil)
	}

	code := c.Query("code")
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	client := conf.Client(context.Background(), token)

	var id string
	var username string
	var email string
	var emailVerified bool
	var connID string

	switch provider {
	case "google":
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var userInfo GoogleOAuth
		if err := sonic.Unmarshal(body, &userInfo); err != nil {
			return err
		}

		id = userInfo.Sub
		username = userInfo.Name
		email = userInfo.Email
		emailVerified = userInfo.EmailVerified
		connID = models.Google.WithID(id)
	}

	if !emailVerified {
		return lib.NewError(fiber.StatusBadRequest, "Email must be verified.", nil)
	}

	var user models.User
	err = lib.DB.Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = models.User{
				ID:        utils.UUIDv4(),
				Username:  username,
				Email:     email,
				MFA:       lib.GenerateString(64),
				Suspended: false,
			}
			if err = lib.DB.Create(&user).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	var conn models.Connection
	err = lib.DB.Where("id = ?", connID).First(&conn).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			conn = models.Connection{
				ID:       connID,
				UserID:   user.ID,
				Verified: true,
			}
			if err = lib.DB.Create(&conn).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	var session = models.Session{
		ID:           lib.GenerateString(64),
		ConnectionID: connID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err = lib.DB.Create(&session).Error; err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    session.ID,
		Path:     "/",
		Domain:   lib.Config.Domain,
		MaxAge:   int(session.ExpiresAt.Sub(time.Now()).Seconds()),
		HTTPOnly: true,
		SameSite: "lax",
	})

	return c.Redirect(fmt.Sprintf("%s/", lib.Config.PublicURL))
}

type GoogleOAuth struct {
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Locale        string `json:"locale"`
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
}
