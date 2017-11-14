package utils

import validator "gopkg.in/go-playground/validator.v9"

// Validator struct that allows a validator to be added to an instance of the echo web framework
type Validator struct {
	validator *validator.Validate
}

// Validate takes a struct and validates it according to the structs validate annotations
func (validator *Validator) Validate(input interface{}) error {
	return validator.validator.Struct(input)
}

// NewValidator creates a Validator instance that can be added to an echo server e.g. echoServer.Validator = NewValidator()
func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}
