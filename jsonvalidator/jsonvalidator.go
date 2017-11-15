package jsonvalidator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	validation "gopkg.in/go-playground/validator.v9"
)

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
	path = strings.Replace(path, "[", ".", 0)
	path = strings.Replace(path, "]", ".", 0)
	fmt.Println("path", path)
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

		return echo.NewHTTPError(http.StatusBadRequest, formattedErrors)
	}

	return nil
}

// NewValidator creates a Validator instance that can be added to an echo server e.g. echoServer.Validator = NewValidator()
func NewValidator() *Validator {
	return &Validator{validator: validation.New()}
}
