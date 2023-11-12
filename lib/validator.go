package lib

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10" // Package for validating struct fields
	"github.com/gofiber/fiber/v2"            // Fiber web framework for Go
)

// Validator instance used to validate struct fields.
var validate = validator.New()

// ParseAndValidate parses the request body into the given struct and performs validation.
func ParseAndValidate(c *fiber.Ctx, body any) error {
	// Parse the body of the request into the provided struct pointer.
	if err := c.BodyParser(body); err != nil {
		return err // Return early with parsing error.
	}

	// Initialize a map to hold field validation errors.
	var fieldErrs = make(map[string][]string)

	// Perform validation on the struct.
	if err := validate.Struct(body); err != nil {
		// Iterate over the field validation errors.
		for _, err := range err.(validator.ValidationErrors) {
			// Retrieve the JSON tag for the field from the struct tag.
			fieldName := jsonTagFromStructField(body, err.StructField())

			// Translate the validation error code to a human-friendly message.
			msg := FieldErrToMsg(err.Tag(), err.Param())

			// Group the messages for each field.
			fieldErrs[fieldName] = append(fieldErrs[fieldName], msg)
		}
	}

	// If there are any validation errors, create and return an error response.
	if len(fieldErrs) > 0 {
		return NewError(http.StatusBadRequest, "Validation errors occurred.", &ErrorDetails{
			Fields: mapToErrorFields(fieldErrs),
		})
	}

	return nil // No errors; validation passed.
}

// jsonTagFromStructField retrieves the JSON tag for a given field in a struct, handling pointer structs.
func jsonTagFromStructField(structPtr any, fieldName string) string {
	// Reflect on the struct to get the field type.
	structValue := reflect.Indirect(reflect.ValueOf(structPtr))
	structType := structValue.Type()

	// Find the field by name and extract its JSON tag.
	if field, found := structType.FieldByName(fieldName); found {
		if tag, ok := field.Tag.Lookup("json"); ok {
			return tag
		}
	}

	// Default to using the field name as is if no JSON tag is found.
	return strings.ToLower(fieldName)
}

// mapToErrorFields converts a map of field errors to a slice of ErrorField structs.
func mapToErrorFields(fieldErrs map[string][]string) []ErrorField {
	var errorFields []ErrorField
	for name, errors := range fieldErrs {
		errorFields = append(errorFields, ErrorField{
			Name:   name,
			Errors: errors,
		})
	}
	return errorFields
}

// FieldErrToMsg converts a validation tag and parameters to a user-friendly error message.
func FieldErrToMsg(tag string, param string) string {
	// Map of custom error messages for each validation tag.
	var tagToMessage = map[string]func(string) string{
		"required": func(_ string) string { return "This field is required." },
		"email":    func(_ string) string { return "This field must contain a valid email address." },
		"min": func(param string) string {
			return fmt.Sprintf("This field must contain at least %s characters.", param)
		},
		// Add more validation tags and their messages as needed.
	}

	// Retrieve the error message function based on the tag, and call it with the parameter.
	if msgFunc, exists := tagToMessage[tag]; exists {
		return msgFunc(param)
	}

	// If no custom message is defined, return the tag name as the error message.
	return tag
}
