package jwt

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// RequiredRoleMiddleware is a echo middleware that will allow to restrict access to a JWT token containing a specific user role
func RequiredRoleMiddleware(publicKey *rsa.PublicKey, requiredRole string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			token := GetTokenFromContext(context)
			parsedToken, err := ParseIfValid(publicKey, token)
			if err != nil {
				return context.JSONBlob(http.StatusUnauthorized, []byte("{\"message\":\"Unauthorized\"}"))
			}

			if !parsedToken.Claims().Has("roles") {
				return context.JSONBlob(http.StatusForbidden, []byte("{\"message\":\"Forbidden\"}"))
			}

			roles := strings.Split(parsedToken.Claims().Get("roles").(string), ",")
			for _, role := range roles {
				if requiredRole == role {
					return next(context)
				}
			}

			return context.JSONBlob(http.StatusForbidden, []byte("{\"message\":\"Forbidden\"}"))
		}
	}
}
