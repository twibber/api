package lib

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func ParseAndValidate(c *fiber.Ctx, body interface{}) (err error) {
	if err = c.BodyParser(&body); err != nil {
		return err
	}

	// multiple errors per field are possible
	var fieldErrs = make(map[string][]string) // key is field, value is an array of all errors
	if err := validate.Struct(body); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			name := strings.ToLower(err.Field())

			// extract struct if pointer
			dataValue := reflect.ValueOf(body)
			if dataValue.Kind() == reflect.Ptr {
				dataValue = dataValue.Elem()
			}

			fieldType, _ := dataValue.Type().FieldByName(err.StructField())

			name = fieldType.Tag.Get("json")

			fieldErrs[name] = append(fieldErrs[name], FieldErrToMsg(err.Tag(), err.Param()))
		}
	}

	// convert map into ErrorField list
	var errs []ErrorField
	for key, element := range fieldErrs {
		errs = append(errs, ErrorField{
			Name:   key,
			Errors: element,
		})
	}

	if len(fieldErrs) > 0 {
		return NewError(http.StatusBadRequest, "The validator has detected errors within the request body.", &ErrorDetails{
			Fields: errs,
		})
	}

	return nil
}

func FieldErrToMsg(err string, param string) string {
	switch err {
	case "required":
		return "This field is required."
	case "email":
		return "This field must contain a valid email address."
	case "min":
		return fmt.Sprintf("This field must contain at least %s characters.", param)
	default:
		return err
	}
}
