package httprequest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/mojlighetsministeriet/utils"
)

// Client extends http.Client by using structs and their validation tags for request/response data
type Client struct {
	http.Client
}

// NewClient creates a http client with timeouts set to 10 seconds and TLS config
func NewClient() (*Client, error) {
	return NewClientWithCustomTimeout(10000)
}

// NewClientWithCustomTimeout creates a http client as NewClient but allows to choosing the timeout in milliseconds
func NewClientWithCustomTimeout(millisecondTimeout time.Duration) (client *Client, err error) {
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

	client = &Client{
		http.Client{
			Timeout:   time.Millisecond * millisecondTimeout,
			Transport: transport,
		},
	}

	return
}

func (client *Client) sendRequest(request *http.Request) (responseBody []byte, err error) {
	response, err := client.Client.Do(request)
	if err != nil {
		return
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		errorBody, _ := ioutil.ReadAll(response.Body)
		err = HTTPError{
			StatusCode:  response.StatusCode,
			ContentType: response.Request.Header.Get("Content-Type"),
			Body:        errorBody,
		}

		return
	}

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(response.Body)
	responseBody = buffer.Bytes()

	return
}

// Get sends a GET request and maps the response to a responseBody struct
func (client *Client) Get(url string) (responseBody []byte, err error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	responseBody, err = client.sendRequest(request)

	return
}

// HTTPError implements an error that retains some additional data about the response
type HTTPError struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

func (err HTTPError) Error() string {
	return strconv.Itoa(err.StatusCode) + " " + http.StatusText(err.StatusCode) + " (" + err.ContentType + "): " + string(err.Body)
}

// JSONClient extends http.Client by using structs and their validation tags for request/response data
type JSONClient struct {
	http.Client
}

func (client *JSONClient) createRequest(method string, url string, requestBody interface{}) (request *http.Request, err error) {
	if requestBody == nil {
		request, err = http.NewRequest(method, url, nil)
	} else {
		requestBodyBuffer := new(bytes.Buffer)
		json.NewEncoder(requestBodyBuffer).Encode(requestBody)
		request, err = http.NewRequest(method, url, requestBodyBuffer)
	}

	return
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

	if responseBody != nil {
		err = json.Unmarshal(buffer, responseBody)
	}

	return
}

// Get sends a GET request and maps the response to a responseBody struct
func (client *JSONClient) Get(url string, responseBody interface{}) (err error) {
	request, err := client.createRequest("GET", url, nil)
	if err != nil {
		return
	}

	err = client.sendRequest(request, responseBody)

	return
}

// Post sends a POST request with a requestBody struct and maps the response to a responseBody struct
func (client *JSONClient) Post(url string, requestBody interface{}, responseBody interface{}) (err error) {
	request, err := client.createRequest("POST", url, requestBody)
	if err != nil {
		return
	}

	err = client.sendRequest(request, responseBody)

	return
}

// Put sends a PUT request with a requestBody struct and maps the response to a responseBody struct
func (client *JSONClient) Put(url string, requestBody interface{}, responseBody interface{}) (err error) {
	request, err := client.createRequest("PUT", url, requestBody)
	if err != nil {
		return
	}

	err = client.sendRequest(request, responseBody)

	return
}

// Delete sends a DELETE request and maps the response to a responseBody struct
func (client *JSONClient) Delete(url string, responseBody interface{}) (err error) {
	request, err := client.createRequest("DELETE", url, nil)
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
