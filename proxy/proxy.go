package proxy

import (
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

func setResponseHeaders(destination, source http.Header) {
	for key, values := range source {
		destination.Del(key)
		for _, value := range values {
			destination.Add(key, value)
		}
	}

	destination.Del("Content-Length")
}

func setProxyRequestHeaders(destination http.Header, sourceContext echo.Context) {
	sourceHeaders := sourceContext.Request().Header
	for key, values := range sourceHeaders {
		destination.Del(key)
		for _, value := range values {
			destination.Add(key, value)
		}
	}

	destination.Add("X-Forwarded-For", sourceContext.RealIP())

	if destination.Get("X-Forwarded-Host") == "" {
		destination.Set("X-Forwarded-Host", sourceContext.Request().Host)
	}

	if destination.Get("X-Forwarded-Proto") == "" {
		destination.Set("X-Forwarded-Proto", sourceContext.Scheme())
	}

	if destination.Get("X-Original-URI") == "" {
		destination.Set("X-Original-URI", sourceContext.Request().URL.RequestURI())
	}
}

func Request(context echo.Context, url string) (err error) {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client := &http.Client{
		Timeout:   time.Second * 10,
		Transport: transport,
	}

	request, err := http.NewRequest(context.Request().Method, url, context.Request().Body)
	if err != nil {
		return
	}
	setProxyRequestHeaders(request.Header, context)
	response, err := client.Do(request)

	if err == nil {
		setResponseHeaders(context.Response().Header(), response.Header)
		return context.Stream(response.StatusCode, response.Header.Get("Content-Type"), response.Body)
	}

	return context.JSONBlob(http.StatusNotFound, []byte("{\"message\":\"Not Found\"}"))
}
