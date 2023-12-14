package email

import (
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/twibber/api/lib"
	"github.com/twibber/api/mailer"
	"github.com/twibber/api/models"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// RegisterDTO defines the structure for registration request data.
type RegisterDTO struct {
	Username string `json:"username"  validate:"required,min=3,max=32,ascii,lowercase"`
	Email    string `json:"email"     validate:"required,email,max=512"`
	Password string `json:"password"  validate:"required,min=8"`
	Captcha  string `json:"captcha"   validate:""`
}

// Register handles the registration of a new user.
func Register(c *fiber.Ctx) error {
	var dto RegisterDTO

	// Parse and validate the request body.
	if err := lib.ParseAndValidate(c, &dto); err != nil {
		return err
	}

	// Check the provided captcha.
	if err := lib.CheckCaptcha(dto.Captcha); err != nil {
		return err
	}

	// Check if the email is already in use.
	var count int64
	if err := lib.DB.Model(models.User{}).Where(models.User{Email: dto.Email}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return lib.ErrEmailExists
	}

	// Hash the password using Argon2id.
	hashedPassword, err := argon2id.CreateHash(dto.Password, &argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  16,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  32,
		KeyLength:   128,
	})
	if err != nil {
		return err
	}

	// Ensure the username is in lowercase.
	dto.Username = strings.ToLower(dto.Username)

	// Check if the username is already taken.
	var usernameCount int64
	if err := lib.DB.Model(models.User{}).Where(models.User{Username: dto.Username}).Count(&usernameCount).Error; err != nil {
		return err
	}
	if usernameCount > 0 {
		return lib.NewError(fiber.StatusConflict, "Username already exists.", &lib.ErrorDetails{
			Fields: []lib.ErrorField{
				{
					Name:   "username",
					Errors: []string{"The username provided is already in use."},
				},
			},
		})
	}

	// Generate a unique token for the user.
	token := lib.GenerateString(64)

	// Generate a TOTP code for email verification.
	totpCode, err := lib.GenerateSecureRandomBase32(32)
	if err != nil {
		return err
	}

	// Start a new database transaction.
	tx := lib.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()

	// Create a new user record.
	user := models.User{
		Username:       dto.Username,
		DisplayName:    dto.Username,
		Admin:          false,
		VerifiedPerson: false,
		Email:          dto.Email,
		MFA:            totpCode,
		Suspended:      false,
	}
	if err := tx.Create(&user).Error; err != nil {
		return err
	}

	// Set expiration for the session.
	exp := 24 * time.Hour

	// Create a new connection record.
	if err := tx.Create(&models.Connection{
		BaseModel: models.BaseModel{ID: models.EmailType.WithID(dto.Email)},
		UserID:    user.ID,
		Password:  hashedPassword,
		Verified:  false,
		Sessions: []models.Session{
			{
				BaseModel: models.BaseModel{ID: token},
				Info: models.SessionInfo{
					IPAddress: c.IP(),
					UserAgent: c.Get("User-Agent"),
				},
				ExpiresAt: time.Now().Add(exp),
			},
		},
	}).Error; err != nil {
		return err
	}

	// Generate a TOTP code for email verification.
	code, err := lib.GenerateTOTP(totpCode, lib.EmailVerification)
	if err != nil {
		return err
	}

	// Send a verification email asynchronously.
	go func() {
		err := mailer.VerifyDTO{
			Defaults: mailer.Defaults{
				Email: user.Email,
				Name:  user.Username,
			},
			Code: code,
		}.Send()
		if err != nil {
			log.WithError(err).Error("new verification code could not be sent")
		}
	}()

	// Set a cookie with the authorization token.
	lib.SetAuth(c, token, exp)

	return c.Status(http.StatusOK).JSON(lib.BlankSuccess)
}
