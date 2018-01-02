package jsonvalidator

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
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

func TestValidateWithEchoVersion(test *testing.T) {
	type User struct {
		Email string `validate:"required,email"`
		Bio   string `validate:"required"`
	}

	structValidator := NewValidatorEcho()

	user := User{Email: "test@example.com", Bio: "I'm a user that has an email."}
	err := structValidator.Validate(user)

	assert.Empty(test, err)
}

func TestValidateDateTime(test *testing.T) {
	type Document struct {
		Title string `validate:"required"`
		Date  string `validate:"required,date-time"`
	}

	structValidator := NewValidator()

	document := Document{Title: "Welcome home!", Date: "2017-03-13T10:22:41Z"}
	err := structValidator.Validate(document)
	assert.NoError(test, err)
}

func TestFailValidateDateTimeWithInvalidFormat(test *testing.T) {
	expectedOutput := `[{"path":"Title", "type":"required"}, {"path":"Date", "type":"date-time"}]`

	type Document struct {
		Title string `validate:"required"`
		Date  string `validate:"required,date-time"`
	}

	structValidator := NewValidator()

	document := Document{Date: "3017-13-13T32:22:41Z"}
	err := structValidator.Validate(document)
	assert.JSONEq(test, expectedOutput, string(err.(ValidationErrors).JSON()))
}

func TestFailValidateOnInvalidEmail(test *testing.T) {
	expectedOutput := `[{"path":"Email", "type": "email"}]`

	type User struct {
		Email string `validate:"required,email"`
		Bio   string `validate:"required"`
	}

	structValidator := NewValidatorEcho()

	user := User{Email: "testexample.com", Bio: "I'm a user that has an email."}
	err := structValidator.Validate(user)
	assert.Error(test, err)

	errEcho, ok := err.(*echo.HTTPError)
	assert.Equal(test, true, ok)
	assert.Equal(test, errEcho.Code, 422)
	assert.JSONEq(test, expectedOutput, string(errEcho.Message.(ValidationErrors).JSON()))
}

func TestFailValidateErrorMethodOnResult(test *testing.T) {
	expectedOutput := "error: email at path Email"

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
	expectedOutput := `[{"path":"wheels.1.radius", "type":"required"}, {"path":"wheels.2.radius", "type":"min"}]`

	type Wheel struct {
		Radius        float64 `json:"radius" validate:"required,min=2"`
		NotSerialized string  `json:"-"`
	}

	type Car struct {
		Brand  string  `json:"brand" validate:"required"`
		Wheels []Wheel `json:"wheels" validate:"required,dive"`
	}

	structValidator := NewValidatorEcho()

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

	errEcho, ok := err.(*echo.HTTPError)
	assert.Equal(test, true, ok)
	assert.Equal(test, errEcho.Code, 422)
	assert.JSONEq(test, expectedOutput, string(errEcho.Message.(ValidationErrors).JSON()))
}

func TestFormatPath(test *testing.T) {
	expectedOutput := "wheels.1.radius"
	output := formatPath("Car.wheels[1].radius")
	assert.Equal(test, expectedOutput, output)
}

func TestNewMalformedJSONResponse(test *testing.T) {
	expectedOutput := `{ "message": "Malformed JSON" }`

	service := echo.New()
	request := httptest.NewRequest(echo.POST, "/", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	context := service.NewContext(request, recorder)
	err := NewMalformedJSONResponse(context)
	assert.NoError(test, err)
	assert.Equal(test, expectedOutput, recorder.Body.String())
	assert.Equal(test, http.StatusBadRequest, recorder.Code)
}
