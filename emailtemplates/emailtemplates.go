package emailtemplates

import (
	"bytes"
	"errors"
	"html/template"
)

// Templates can hold templates
type Templates struct {
	templates []*Template
}

// Add adds a new template
func (templates *Templates) Add(emailTemplate Template) {
	emailTemplate.Compile()
	templates.templates = append(templates.templates, &emailTemplate)
}

// Render will generate and return the subject and body strings from a template
func (templates *Templates) Render(name string, to string, subjectData interface{}, bodyData interface{}) (result Email, err error) {
	var templateToRender *Template

	for _, emailTemplate := range templates.templates {
		if emailTemplate.Name == name {
			templateToRender = emailTemplate
			break
		}
	}

	if templateToRender == nil {
		err = errors.New("Template " + name + " is not registered")
		return
	}

	subject, err := templateToRender.GetSubject(subjectData)
	if err != nil {
		subject = ""
		return
	}

	body, err := templateToRender.GetBody(bodyData)
	if err != nil {
		subject = ""
		body = ""
		return
	}

	result = Email{To: to, Subject: subject, Body: body}

	return
}

// Template lets you format emails with go templates
type Template struct {
	Name            string
	Subject         string
	subjectTemplate *template.Template
	Body            string
	bodyTemplate    *template.Template
}

// Email can hold the result of a rendered template
type Email struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
	To      string `json:"to"`
}

// Compile will re-compile the template Subject and Body to templates
func (emailTemplate *Template) Compile() (err error) {
	emailTemplate.subjectTemplate, err = template.New("subject").Parse(emailTemplate.Subject)
	if err != nil {
		return
	}

	emailTemplate.bodyTemplate, err = template.New("body").Parse(emailTemplate.Body)
	return
}

func (emailTemplate *Template) execute(subTemplate *template.Template, data interface{}) (output string, err error) {
	buffer := new(bytes.Buffer)

	if subTemplate == nil {
		err = emailTemplate.Compile()
		if err != nil {
			return
		}
	}

	err = subTemplate.Execute(buffer, data)
	if err != nil {
		return
	}

	output = buffer.String()

	return
}

// GetSubject will return the populated Subject by passing a data map to the function
func (emailTemplate *Template) GetSubject(data interface{}) (output string, err error) {
	output, err = emailTemplate.execute(emailTemplate.subjectTemplate, data)
	return
}

// GetBody will return the populated Body by passing a data map to the function
func (emailTemplate *Template) GetBody(data interface{}) (output string, err error) {
	output, err = emailTemplate.execute(emailTemplate.bodyTemplate, data)
	return
}
