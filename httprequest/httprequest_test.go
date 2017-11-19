package httprequest

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestJSONClientGet(test *testing.T) {
	url := "https://api.github.com/repos/mojlighetsministeriet/utils"
	type Response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	client, err := NewJSONClient()
	assert.NoError(test, err)

	response := Response{}
	err = client.Get(url, &response)
	assert.NoError(test, err)
	assert.Equal(test, 103898330, response.ID)
	assert.Equal(test, "utils", response.Name)
}

func TestFailJSONClientGetWithNotFoundURL(test *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://mojlighetsministeriet.se",
		httpmock.NewStringResponder(404, `{"message": "Not Found"}`))

	url := "https://mojlighetsministeriet.se"

	client, err := NewJSONClient()
	assert.NoError(test, err)

	err = client.Get(url, nil)

	assert.Error(test, err)
	assert.Equal(test, "404 Not Found (application/json; charset=utf-8): {\"message\":\"Not Found\",\"documentation_url\":\"https://developer.github.com/v3\"}", err.Error())
}

func TestFailJSONClientGetWithBadURL(test *testing.T) {
	url := "http://this¤{I}sawierdstring"
	type Response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	client, err := NewJSONClient()
	assert.NoError(test, err)

	response := Response{}
	err = client.Get(url, &response)

	assert.Error(test, err)
	assert.Equal(test, "parse http://this¤{I}sawierdstring: invalid character \"{\" in host name", err.Error())
	assert.Equal(test, 0, response.ID)
	assert.Equal(test, "", response.Name)
}

func TestJSONClientDelete(test *testing.T) {
	url := "https://api.github.com/repos/mojlighetsministeriet/utils"
	type Response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	client, err := NewJSONClient()
	assert.NoError(test, err)

	response := Response{}
	err = client.Get(url, &response)
	assert.NoError(test, err)
	assert.Equal(test, 103898330, response.ID)
	assert.Equal(test, "utils", response.Name)
}

func TestFailJSONClientDeleteWithBadURL(test *testing.T) {
	url := "http://this¤{I}sawierdstring"
	type Response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	client, err := NewJSONClient()
	assert.NoError(test, err)

	response := Response{}
	err = client.Get(url, &response)

	assert.Error(test, err)
	assert.Equal(test, "parse http://this¤{I}sawierdstring: invalid character \"{\" in host name", err.Error())
	assert.Equal(test, 0, response.ID)
	assert.Equal(test, "", response.Name)
}

func TestFailJSONClientGetWithInvalidDomainName(test *testing.T) {
	url := "https://adomainthatdoesnotexist"
	type Response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	client, err := NewJSONClient()
	assert.NoError(test, err)

	response := Response{}
	err = client.Get(url, &response)

	assert.Error(test, err)
	assert.Equal(test, "Get https://adomainthatdoesnotexist: dial tcp: lookup adomainthatdoesnotexist: no such host", err.Error())
	assert.Equal(test, 0, response.ID)
	assert.Equal(test, "", response.Name)
}
