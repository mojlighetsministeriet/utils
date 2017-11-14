package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(test *testing.T) {
	type User struct {
		Email string `validate:"required,email"`
		Bio   string `validate:"required"`
	}

	structValidator := NewValidator()

	user := User{Email: "test@example.com", Bio: "I'm a user that has an email."}
	err := structValidator.Validate(user)
	assert.NoError(test, err)
}

func TestFailValidateOnInvalidEmail(test *testing.T) {
	expectedOutput := "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"

	type User struct {
		Email string `validate:"required,email"`
		Bio   string `validate:"required"`
	}

	structValidator := NewValidator()

	user := User{Email: "testexample.com", Bio: "I'm a user that has an email."}
	err := structValidator.Validate(user)
	assert.Error(test, err)
	assert.Equal(test, expectedOutput, err.Error())
}
