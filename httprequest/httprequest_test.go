package httprequest

import (
	"net/http"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func startTestServer() *echo.Echo {
	server := echo.New()
	server.HideBanner = true

	server.GET("/", func(context echo.Context) error {
		return context.JSONBlob(http.StatusOK, []byte("{\"id\":103898330,\"name\":\"utils\"}"))
	})

	server.POST("/", func(context echo.Context) error {
		return context.JSONBlob(http.StatusCreated, []byte("{\"id\":103898330,\"name\":\"utils\"}"))
	})

	server.PUT("/iexist", func(context echo.Context) error {
		return context.JSONBlob(http.StatusCreated, []byte("{\"id\":103898330,\"name\":\"utils\"}"))
	})

	server.DELETE("/deleteme", func(context echo.Context) error {
		return context.JSONBlob(http.StatusOK, []byte("{\"message\":\"Deleted\"}"))
	})

	go func(server *echo.Echo) {
		server.Start(":12345")
	}(server)

	return server
}

func TestJSONClientGet(test *testing.T) {
	server := startTestServer()
	defer server.Close()

	url := "http://localhost:12345"
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
	server := startTestServer()
	defer server.Close()

	url := "http://localhost:12345/idonotexist"

	client, err := NewJSONClient()
	assert.NoError(test, err)

	err = client.Get(url, nil)

	assert.Error(test, err)
	assert.Equal(test, "404 Not Found (application/json; charset=utf-8): {\"message\":\"Not Found\"}", err.Error())
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
	server := startTestServer()
	defer server.Close()

	url := "http://localhost:12345/deleteme"

	client, err := NewJSONClient()
	assert.NoError(test, err)

	err = client.Delete(url, nil)
	assert.NoError(test, err)
}

func TestFailJSONClientDeleteWithBadURL(test *testing.T) {
	url := "http://this¤{I}sawierdstring"

	client, err := NewJSONClient()
	assert.NoError(test, err)

	err = client.Delete(url, nil)

	assert.Error(test, err)
	assert.Equal(test, "parse http://this¤{I}sawierdstring: invalid character \"{\" in host name", err.Error())
}

func TestJSONClientPost(test *testing.T) {
	server := startTestServer()
	defer server.Close()

	url := "http://localhost:12345/"

	client, err := NewJSONClient()
	assert.NoError(test, err)

	type Request struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	err = client.Post(url, Request{ID: 12345, Name: "A name"}, nil)
	assert.NoError(test, err)
}

func TestFailJSONClientPostWithBadURL(test *testing.T) {
	url := "http://this¤{I}sawierdstring"

	client, err := NewJSONClient()
	assert.NoError(test, err)

	err = client.Post(url, nil, nil)

	assert.Error(test, err)
	assert.Equal(test, "parse http://this¤{I}sawierdstring: invalid character \"{\" in host name", err.Error())
}

func TestJSONClientPut(test *testing.T) {
	server := startTestServer()
	defer server.Close()

	url := "http://localhost:12345/iexist"

	client, err := NewJSONClient()
	assert.NoError(test, err)

	err = client.Put(url, nil, nil)
	assert.NoError(test, err)
}

func TestFailJSONClientPutWithBadURL(test *testing.T) {
	url := "http://this¤{I}sawierdstring"

	client, err := NewJSONClient()
	assert.NoError(test, err)

	err = client.Put(url, nil, nil)

	assert.Error(test, err)
	assert.Equal(test, "parse http://this¤{I}sawierdstring: invalid character \"{\" in host name", err.Error())
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
	assert.Equal(test, "Get https://adomainthatdoesnotexist: dial tcp: lookup adomainthatdoesnotexist: No address associated with hostname", err.Error())
	assert.Equal(test, 0, response.ID)
	assert.Equal(test, "", response.Name)
}
