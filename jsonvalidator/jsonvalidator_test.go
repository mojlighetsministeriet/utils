package jsonvalidator

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
	expectedOutput := "code=400, message=[{Email email}]"

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

func TestFailValidateOnNestedStructure(test *testing.T) {
	expectedOutput := "code=400, message=[{wheels.1.radius required} {wheels.2.radius min}]"

	type Wheel struct {
		Radius        float64 `json:"radius" validate:"required,min=2"`
		NotSerialized string  `json:"-"`
	}

	type Car struct {
		Brand  string  `json:"brand" validate:"required"`
		Wheels []Wheel `json:"wheels" validate:"required,dive"`
	}

	structValidator := NewValidator()

	car := Car{
		Brand: "tesla",
		Wheels: []Wheel{
			Wheel{Radius: 5},
			Wheel{},
			Wheel{Radius: 1},
		},
	}
	err := structValidator.Validate(car)
	assert.Error(test, err)
	assert.Equal(test, expectedOutput, err.Error())
}

func TestFormatPath(test *testing.T) {
	expectedOutput := "wheels.1.radius"
	output := formatPath("Car.wheels[1].radius")
	assert.Equal(test, expectedOutput, output)
}
