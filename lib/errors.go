package lib

// this is not a separate package to avoid import cycle (I hate my life)

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/bytedance/sonic/decoder"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// these errors are only the recurring ones, if it's an error that's only used once it's just created with NewError on the return
var (
	ErrInternal           = NewError(http.StatusNotImplemented, "An internal server error occurred while attempting to process the request.", nil)
	ErrForbidden          = NewError(http.StatusForbidden, "You do not have permission to access the requested resource.", nil)
	ErrUnauthorised       = NewError(http.StatusUnauthorized, "You are not authorised to access this endpoint.", nil)
	ErrNotFound           = NewError(http.StatusNotFound, "The requested resource does not exist.", nil)
	ErrNotImplemented     = NewError(http.StatusNotImplemented, "A portion of this request has not been implemented.", nil)
	ErrInvalidCredentials = NewError(http.StatusUnauthorized, "Invalid credentials. Please try again.", &ErrorDetails{
		Fields: []ErrorField{
			{
				Name:   "email",
				Errors: []string{"Invalid credentials. Please try again."},
			},
			{
				Name:   "password",
				Errors: []string{"Invalid credentials. Please try again."},
			},
		},
	})
	ErrInvalidCaptcha = NewError(http.StatusUnauthorized, "The captcha response suggests this action was not performed by a human.", &ErrorDetails{
		Fields: []ErrorField{
			{
				Name:   "captcha",
				Errors: []string{"The captcha response suggests this action was not performed by a human."},
			},
		},
	})
	ErrEmailExists = NewError(http.StatusConflict, "The email address provided has already been registered.", &ErrorDetails{
		Fields: []ErrorField{
			{
				Name:   "email",
				Errors: []string{"The email address provided has already been registered."},
			},
		},
	})
)

type Error struct {
	Status  int           `json:"-"`
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details *ErrorDetails `json:"details,omitempty"`
}

func (e Error) Error() string {
	return fmt.Sprintf("Code: %s, Message: %s, Details: %v", e.Code, e.Message, e.Details)
}

type ErrorDetails struct {
	Fields []ErrorField `json:"fields,omitempty"`
	Debug  any          `json:"debug,omitempty"`
}

type ErrorField struct {
	Name   string   `json:"name"`
	Errors []string `json:"errors"`
}

func NewError(status int, message string, details *ErrorDetails, code ...string) Error {
	var statusCode string
	if len(code) > 0 {
		statusCode = code[0]
	} else {
		statusCode = strings.ReplaceAll(strings.ToUpper(utils.StatusMessage(status)), " ", "_")
	}

	return Error{
		Status:  status,
		Code:    statusCode,
		Message: message,
		Details: details,
	}
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	switch err.(type) {
	case Error:
		e := err.(Error)

		return c.Status(e.Status).JSON(Response{
			Success:    false,
			ObjectName: "error",
			Data:       e,
		})
	case *fiber.Error:
		fiberErr := err.(*fiber.Error)
		e := NewError(fiberErr.Code, fiberErr.Message, nil)

		return c.Status(e.Status).JSON(Response{
			Success:    false,
			ObjectName: "error",
			Data:       e,
		})
	case decoder.SyntaxError:
		e := NewError(fiber.StatusUnprocessableEntity, "Invalid JSON was provided in the request body.", nil)

		return c.Status(e.Status).JSON(Response{
			Success:    false,
			ObjectName: "error",
			Data:       e,
		})
	default:
		switch err {
		case gorm.ErrRecordNotFound, gorm.ErrEmptySlice:
			var e Error
			switch c.Route().Path {
			case "/auth/email/register", "/auth/email/login":
				e = ErrInvalidCredentials
			default:
				e = ErrNotFound
			}

			return c.Status(e.Status).JSON(Response{
				Success:    false,
				ObjectName: "error",
				Data:       e,
			})
		default:
			// this is an unhandled error
			// only give true error while debugging, to avoid leaking sensitive information
			e := NewError(fiber.StatusInternalServerError, "An internal server error occurred while processing your request", nil)
			if Config.Debug {
				log.WithError(err).WithField("errType", reflect.TypeOf(err).String()).Error("An unhandled error occurred.")
				e = NewError(fiber.StatusInternalServerError, fmt.Sprintf("%s: %s", reflect.TypeOf(err).String(), err.Error()), nil)
			}

			return c.Status(e.Status).JSON(Response{
				Success:    false,
				ObjectName: "error",
				Data:       e,
			})
		}
	}
}
