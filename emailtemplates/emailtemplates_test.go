package emailtemplates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplatesRender(test *testing.T) {
	templates := Templates{}
	template := Template{
		Name:    "newAccount",
		Subject: "Hi {{.Name}}! this is your new account.",
		Body:    "You have a new account, choose your password <a href=\"{{.ServiceURL}}/api/reset-password/{{.ResetToken}}\" target=\"_blank\">here</a>.",
	}
	templates.Add(template)

	data := make(map[string]string)
	data["Name"] = "Anna"
	data["ServiceURL"] = "https://internt.mojlighetsministeriet.se"
	data["ResetToken"] = "e980ad6b-3b78-5579-804a-a3ec18798332"

	result, err := templates.Render("newAccount", data, data)

	assert.NoError(test, err)
	assert.Equal(test, "Hi Anna! this is your new account.", result.Subject)
	assert.Equal(test, "You have a new account, choose your password <a href=\"https://internt.mojlighetsministeriet.se/api/reset-password/e980ad6b-3b78-5579-804a-a3ec18798332\" target=\"_blank\">here</a>.", result.Body)
}

func TestFailTemplatesRenderWithMissingTemplate(test *testing.T) {
	templates := Templates{}

	result, err := templates.Render("missingTemplate", nil, nil)

	assert.Error(test, err)
	assert.Equal(test, "Template missingTemplate is not registered", err.Error())
	assert.Equal(test, "", result.Subject)
	assert.Equal(test, "", result.Body)
}

func TestFailTemplatesRenderWithInvalidSubjectTemplate(test *testing.T) {
	templates := Templates{}
	template := Template{
		Name:    "newAccount",
		Subject: "Hi {{.Na{me}}! this is your new account.",
		Body:    "You have a new account, choose your password <a href=\"{{.ServiceURL}}/api/reset-password/{{.ResetToken}}\" target=\"_blank\">here</a>.",
	}
	templates.Add(template)

	data := make(map[string]string)
	data["Name"] = "Anna"
	data["ServiceURL"] = "https://internt.mojlighetsministeriet.se"
	data["ResetToken"] = "e980ad6b-3b78-5579-804a-a3ec18798332"

	result, err := templates.Render("newAccount", data, data)

	assert.Error(test, err)
	assert.Equal(test, "template: subject:1: unexpected bad character U+007B '{' in command", err.Error())
	assert.Equal(test, "", result.Subject)
	assert.Equal(test, "", result.Body)
}

func TestFailTemplatesRenderWithInvalidBodyTemplate(test *testing.T) {
	templates := Templates{}
	template := Template{
		Name:    "newAccount",
		Subject: "Hi {{.Name}}! this is your new account.",
		Body:    "You have a new account, choose your password <a href=\"{{.Servi{ceURL}}/api/reset-password/{{.ResetToken}}\" target=\"_blank\">here</a>.",
	}
	templates.Add(template)

	data := make(map[string]string)
	data["Name"] = "Anna"
	data["ServiceURL"] = "https://internt.mojlighetsministeriet.se"
	data["ResetToken"] = "e980ad6b-3b78-5579-804a-a3ec18798332"

	result, err := templates.Render("newAccount", data, data)

	assert.Error(test, err)
	assert.Equal(test, "template: body:1: unexpected bad character U+007B '{' in command", err.Error())
	assert.Equal(test, "", result.Subject)
	assert.Equal(test, "", result.Body)
}
