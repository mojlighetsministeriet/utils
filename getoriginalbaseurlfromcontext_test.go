package utils

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetOriginalSystemURLFromContext(test *testing.T) {
	service := echo.New()
	request := httptest.NewRequest(echo.GET, "/", strings.NewReader("{}"))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Set("X-Forwarded-Proto", "https")
	request.Header.Set("X-Forwarded-Host", "internt.mojlighetsministeriet.se")
	recorder := httptest.NewRecorder()
	context := service.NewContext(request, recorder)

	url := GetOriginalSystemURLFromContext(context)

	assert.Equal(test, "https://internt.mojlighetsministeriet.se", url)
}
