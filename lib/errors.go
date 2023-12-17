package lib

import (
	"encoding/json"
	"fmt"
	cfg "github.com/twibber/api/config"
	"net/http"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"       // Fiber web framework
	"github.com/gofiber/fiber/v2/utils" // Utility functions for Fiber
	log "github.com/sirupsen/logrus"    // Structured logging package
	"gorm.io/gorm"                      // GORM ORM package
)

// Predefined errors for common API responses
var (
	ErrInternal           = NewError(http.StatusNotImplemented, "An internal server error occurred while attempting to process the request.", nil)
	ErrForbidden          = NewError(http.StatusForbidden, "You do not have permission to access the requested resource.", nil)
	ErrUnauthorised       = NewError(http.StatusUnauthorized, "You are not authorised to access this endpoint.", nil)
	ErrNotFound           = NewError(http.StatusNotFound, "The requested resource does not exist.", nil)
	ErrNotImplemented     = NewError(http.StatusNotImplemented, "A portion of this request has not been implemented.", nil)
	ErrInvalidCredentials = NewError(http.StatusBadRequest, "Invalid credentials. Please try again.", &ErrorDetails{
		Fields: []ErrorField{
			{Name: "email", Errors: []string{"Invalid credentials. Please try again."}},
			{Name: "password", Errors: []string{"Invalid credentials. Please try again."}},
		},
	})
	ErrInvalidCaptcha = NewError(http.StatusBadRequest, "The captcha response suggests this action was not performed by a human.", &ErrorDetails{
		Fields: []ErrorField{
			{Name: "captcha", Errors: []string{"The captcha response suggests this action was not performed by a human."}},
		},
	})
	ErrEmailExists = NewError(http.StatusConflict, "The email address provided has already been registered.", &ErrorDetails{
		Fields: []ErrorField{
			{Name: "email", Errors: []string{"The email address provided has already been registered."}},
		},
	})
)

// Error represents a standardised error response for the API.
type Error struct {
	Status  int           `json:"-"`                 // HTTP status code, not included in the response
	Code    string        `json:"code"`              // API-specific error code
	Message string        `json:"message"`           // Human-readable error message
	Details *ErrorDetails `json:"details,omitempty"` // Optional details about the error
}

// Error formats the error message string.
func (e Error) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s, Details: %v", e.Code, e.Message, e.Details)
}

// ErrorDetails holds additional data about the error.
type ErrorDetails struct {
	Fields []ErrorField `json:"fields,omitempty"` // Specific fields related to the error
	Debug  any          `json:"debug,omitempty"`  // Debug information, included only if debugging is enabled
}

// ErrorField provides detailed errors for specific fields in the request.
type ErrorField struct {
	Name   string   `json:"name"`   // Name of the field
	Errors []string `json:"errors"` // List of error messages for the field
}

// NewError creates a new Error with the provided status, message, and optional details.
func NewError(status int, message string, details *ErrorDetails, code ...string) Error {
	var statusCode string
	if len(code) > 0 {
		statusCode = code[0] // Use provided error code if available
	} else {
		// Otherwise, generate an error code from the HTTP status message
		statusCode = strings.ReplaceAll(strings.ToUpper(utils.StatusMessage(status)), " ", "_")
	}

	return Error{
		Status:  status,
		Code:    statusCode,
		Message: message,
		Details: details,
	}
}

// ErrorHandler is a custom error handler for the Fiber application.
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Handles different types of errors and formats them for API responses
	switch err.(type) {
	case Error:
		e := err.(Error)
		return c.Status(e.Status).JSON(Response{Success: false, Data: e})
	case *fiber.Error:
		fiberErr := err.(*fiber.Error)
		e := NewError(fiberErr.Code, fiberErr.Message, nil)
		return c.Status(e.Status).JSON(Response{Success: false, Data: e})
	case *json.SyntaxError:
		var e Error
		if cfg.Config.Debug {
			e = NewError(fiber.StatusBadRequest, fmt.Sprintf("%s: %s", reflect.TypeOf(err).String(), err.Error()), nil)
		} else {
			e = NewError(fiber.StatusUnprocessableEntity, "Invalid JSON was provided in the request body.", nil)
		}

		return c.Status(e.Status).JSON(Response{Success: false, Data: e})
	default:
		// Handles GORM's specific errors and maps them to API errors
		switch err {
		case gorm.ErrRecordNotFound, gorm.ErrEmptySlice:
			var e Error
			switch c.Route().Path {
			case "/auth/email/register", "/auth/email/login":
				e = ErrInvalidCredentials
			default:
				e = ErrNotFound
			}
			return c.Status(e.Status).JSON(Response{Success: false, Data: e})
		default:
			// Logs unhandled errors and returns a generic error response
			e := NewError(fiber.StatusInternalServerError, "An internal server error occurred while processing your request", nil)
			if cfg.Config.Debug {
				var debugInfo any

				debugInfo = log.Fields{
					"ErrorType": reflect.TypeOf(err).String(),
					"Message":   err.Error(),
				}

				log.WithError(err).WithFields(log.Fields{
					"errType": reflect.TypeOf(err).String(),
					"error":   err.Error(),
				}).Error("An unhandled error occurred.")

				e = NewError(fiber.StatusInternalServerError, "An internal server error occurred while processing your request", &ErrorDetails{
					Debug: debugInfo,
				})
			}
			return c.Status(e.Status).JSON(Response{Success: false, Data: e})
		}
	}
}
