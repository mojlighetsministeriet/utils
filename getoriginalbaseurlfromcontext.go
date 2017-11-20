package utils

import "github.com/labstack/echo"

// GetOriginalSystemURLFromContext will take a context and return the systems base URL (the URL to the actual external host) e.g. https://internt.mojlighetsministeriet.se
func GetOriginalSystemURLFromContext(context echo.Context) string {
	headers := context.Request().Header
	return headers.Get("X-Forwarded-Proto") + "://" + headers.Get("X-Forwarded-Host")
}
