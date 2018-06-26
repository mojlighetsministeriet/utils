package server

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo"
)

func RemoveExtraSlashesMiddleware() echo.MiddlewareFunc {
	slashPattern := regexp.MustCompile("//+")

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			request := context.Request()
			url := request.URL
			path := url.Path
			queryString := context.QueryString()
			cleanedPath := strings.TrimSuffix(slashPattern.ReplaceAllString(path, "/"), "/")

			if len(cleanedPath) > 0 && cleanedPath != "/" && cleanedPath != path {
				uri := cleanedPath

				if queryString != "" {
					uri += "?" + queryString
				}

				return context.Redirect(http.StatusMovedPermanently, uri)
			}

			return next(context)
		}
	}
}
