package jsonvalidator

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/labstack/echo"
	validation "gopkg.in/go-playground/validator.v9"
)

// TODO: Add extra properties for some validation errors such as min, max, min length, max length explaining what the limits are

// ValidationError is the JSON representation of a validation error
type ValidationError struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// Validator struct that allows a validator to be added to an instance of the echo web framework
type Validator struct {
	validator *validation.Validate
}

func formatPath(path string) string {
	path = strings.Replace(path, "[", ".", -1)
	path = strings.Replace(path, "].", ".", -1)
	stripFirstPathNode := regexp.MustCompile(`^\s*\w+\.`)
	path = stripFirstPathNode.ReplaceAllString(path, "")
	return path
}

// Validate takes a struct and validates it according to the structs validate annotations
func (validator *Validator) Validate(input interface{}) error {
	validationError := validator.validator.Struct(input)
	if validationError != nil {
		formattedErrors := []ValidationError{}
		fieldErrors := validationError.(validation.ValidationErrors)
		for _, fieldError := range fieldErrors {
			formattedErrors = append(formattedErrors, ValidationError{
				Path: formatPath(fieldError.Namespace()),
				Type: fieldError.Tag(),
			})
		}

		return echo.NewHTTPError(422, formattedErrors)
	}

	return nil
}

// NewValidator creates a Validator instance that can be added to an echo server e.g. echoServer.Validator = NewValidator()
func NewValidator() *Validator {
	validator := &Validator{validator: validation.New()}
	validator.validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return validator
}

// NewMalformedJSONResponse will create a JSON response that indicates a client has sent a JSON body that has syntax errors
func NewMalformedJSONResponse(context echo.Context) error {
	return context.JSONBlob(http.StatusBadRequest, []byte(`{ "message": "Malformed JSON" }`))
}
