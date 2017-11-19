package httprequest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/mojlighetsministeriet/utils"
)

// NewClient creates a http client with timeouts set to 10 seconds and TLS config
func NewClient() (*http.Client, error) {
	return NewClientWithCustomTimeout(10000)
}

// NewClientWithCustomTimeout creates a http client as NewClient but allows to choosing the timeout in milliseconds
func NewClientWithCustomTimeout(millisecondTimeout time.Duration) (client *http.Client, err error) {
	tlsConfig, err := utils.GetCACertificatesTLSConfig()
	if err != nil {
		return
	}

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     tlsConfig,
	}

	client = &http.Client{
		Timeout:   time.Millisecond * millisecondTimeout,
		Transport: transport,
	}

	return
}

// HTTPError implements an error that retains some additional data about the response
type HTTPError struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

func (err HTTPError) Error() string {
	return string(err.StatusCode) + " " + http.StatusText(err.StatusCode) + " (" + err.ContentType + "): " + string(err.Body)
}

// JSONClient extends http.Client by using structs and their validation tags for request/response data
type JSONClient struct {
	http.Client
}

func (client *JSONClient) sendRequest(request *http.Request, responseBody interface{}) (err error) {
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := client.Client.Do(request)
	if err != nil {
		return
	}

	buffer, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		err = HTTPError{
			StatusCode:  response.StatusCode,
			ContentType: response.Request.Header.Get("Content-Type"),
			Body:        buffer,
		}

		return
	}

	err = json.Unmarshal(buffer, responseBody)

	return
}

// Get sends a GET request and maps the response to a responseBody struct
func (client *JSONClient) Get(url string, responseBody interface{}) (err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	err = client.sendRequest(request, responseBody)

	return
}

// Post sends a POST request with a requestBody struct and maps the response to a responseBody struct
func (client *JSONClient) Post(url string, requestBody interface{}, responseBody interface{}) (err error) {
	requestBodyBuffer := new(bytes.Buffer)
	json.NewEncoder(requestBodyBuffer).Encode(requestBody)

	request, err := http.NewRequest("POST", url, requestBodyBuffer)
	if err != nil {
		return
	}

	err = client.sendRequest(request, responseBody)

	return
}

// Put sends a PUT request with a requestBody struct and maps the response to a responseBody struct
func (client *JSONClient) Put(url string, requestBody interface{}, responseBody interface{}) (err error) {
	requestBodyBuffer := new(bytes.Buffer)
	json.NewEncoder(requestBodyBuffer).Encode(requestBody)

	request, err := http.NewRequest("PUT", url, requestBodyBuffer)
	if err != nil {
		return
	}

	err = client.sendRequest(request, responseBody)

	return
}

// Delete sends a DELETE request and maps the response to a responseBody struct
func (client *JSONClient) Delete(url string, responseBody interface{}) (err error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}

	err = client.sendRequest(request, responseBody)

	return
}

// NewJSONClient creates a http client for JSON requests with timeouts set to 10 seconds and TLS config
func NewJSONClient() (*JSONClient, error) {
	return NewJSONClientWithCustomTimeout(10000)
}

// NewJSONClientWithCustomTimeout creates a http client for JSON requests as NewJSONClient but allows to choosing the timeout in milliseconds
func NewJSONClientWithCustomTimeout(millisecondTimeout time.Duration) (client *JSONClient, err error) {
	tlsConfig, err := utils.GetCACertificatesTLSConfig()
	if err != nil {
		return
	}

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		TLSClientConfig:     tlsConfig,
	}

	client = &JSONClient{
		http.Client{
			Timeout:   time.Millisecond * millisecondTimeout,
			Transport: transport,
		},
	}

	return
}
