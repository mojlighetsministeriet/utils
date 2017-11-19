package httprequest

import (
	"testing"

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

func TestFailJSONClientGetWithInvalidURL(test *testing.T) {
	url := "https://api.github.com/repos/mojlighetsministeriet/notfoundurl"
	type Response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	client, err := NewJSONClient()
	assert.NoError(test, err)

	response := Response{}
	err = client.Get(url, &response)

	assert.Error(test, err)
	assert.Equal(test, "", err.Error())
	assert.Equal(test, 0, response.ID)
	assert.Equal(test, "", response.Name)
}

func TestFailJSONClientGetWithoutProtocol(test *testing.T) {
	url := "this¤{I}sawierdstring"
	type Response struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	client, err := NewJSONClient()
	assert.NoError(test, err)

	response := Response{}
	err = client.Get(url, &response)

	assert.NoError(test, err)
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
