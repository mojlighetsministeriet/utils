package jsonvalidator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo"
	validation "gopkg.in/go-playground/validator.v9"
)

// TODO: Add extra properties for some validation errors such as min, max, min length, max length explaining what the limits are

// ValidationError is the JSON representation of a validation error
type ValidationError struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// Error is used for development debugging of validation errors
func (err ValidationError) Error() string {
	return "error: " + err.Type + " at path " + err.Path
}

// ValidationErrors is an error containing several ValidationError instances
type ValidationErrors []ValidationError

// Error is intended only to be used for development.
// To get all errors in a structured way, use errors := err.(ValidationErrors)
func (validationErrors ValidationErrors) Error() string {
	buffer := bytes.NewBufferString("")

	for i := 0; i < len(validationErrors); i++ {
		buffer.WriteString(validationErrors[i].Error())
		buffer.WriteString("\n")
	}

	return strings.TrimSpace(buffer.String())
}

// JSON encodes the errors to a JSON []byte slice
func (validationErrors ValidationErrors) JSON() []byte {
	data, _ := json.Marshal(validationErrors)
	return data
}

// Validator struct that simplifies validating structs
type Validator struct {
	validator *validation.Validate
}

// Validate takes a struct and validates it according to the structs validate annotations
func (validator *Validator) Validate(input interface{}) error {
	validationError := validator.validator.Struct(input)
	if validationError != nil {
		formattedErrors := ValidationErrors{}
		fieldErrors := validationError.(validation.ValidationErrors)
		for _, fieldError := range fieldErrors {
			formattedErrors = append(formattedErrors, ValidationError{
				Path: formatPath(fieldError.Namespace()),
				Type: fieldError.Tag(),
			})
		}

		return formattedErrors
	}

	return nil
}

// NewValidator creates a Validator instance
func NewValidator() *Validator {
	validator := &Validator{validator: validation.New()}

	validator.validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		// TODO: make this check more clear
		if name == "-" {
			return ""
		}
		return name
	})

	validator.validator.RegisterValidation("date-time", func(fieldLevel validation.FieldLevel) bool {
		value := fieldLevel.Field().String()
		_, err := time.Parse(time.RFC3339Nano, value)

		if err == nil {
			return true
		}
		return false
	})

	return validator
}

// ValidatorEcho struct that allows a validator to be added to an instance of the echo web framework
type ValidatorEcho struct {
	*Validator
}

// Validate takes a struct and validates it according to the structs validate annotations
func (validator *ValidatorEcho) Validate(input interface{}) (err *echo.HTTPError) {
	validationError := validator.Validator.Validate(input)
	if validationError == nil {
		return
	}

	errors := validationError.(ValidationErrors)
	if errors != nil {
		return echo.NewHTTPError(422, errors)
	}

	return nil
}

// NewValidatorEcho creates a Validator instance that can be added to an echo server e.g. echoServer.Validator = NewValidator()
func NewValidatorEcho() *ValidatorEcho {
	validator := &Validator{validator: validation.New()}
	validator.validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &ValidatorEcho{Validator: validator}
}

// NewMalformedJSONResponse will create a JSON response that indicates a client has sent a JSON body that has syntax errors
func NewMalformedJSONResponse(context echo.Context) error {
	return context.JSONBlob(http.StatusBadRequest, []byte(`{ "message": "Malformed JSON" }`))
}

func formatPath(path string) string {
	path = strings.Replace(path, "[", ".", -1)
	path = strings.Replace(path, "].", ".", -1)
	stripFirstPathNode := regexp.MustCompile(`^\s*\w+\.`)
	path = stripFirstPathNode.ReplaceAllString(path, "")
	return path
}
