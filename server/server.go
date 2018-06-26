package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	*echo.Echo
	useTLS bool
}

func (server *Server) UseTLS() bool {
	return server.useTLS
}

func (server *Server) Listen(address string) {
	server.addHelpResource()

	if server.useTLS {
		server.Logger.Fatal(server.StartAutoTLS(address))
	} else {
		server.Logger.Fatal(server.Start(address))
	}
}

func (server *Server) addHelpResource() {
	type routeInfo struct {
		Path   string `json:"path"`
		Method string `json:"method"`
	}

	var registeredRoutes []routeInfo

	for _, route := range server.Routes() {
		if !strings.HasSuffix(route.Path, "/*") {
			registeredRoute := routeInfo{
				Path:   route.Path,
				Method: route.Method,
			}
			registeredRoutes = append(registeredRoutes, registeredRoute)
		}
	}

	server.GET("/help", func(context echo.Context) error {
		return context.JSON(http.StatusOK, registeredRoutes)
	})
}

func NewServer(useTLS bool, behindProxy bool, bodyLimit string) *Server {
	server := Server{Echo: echo.New()}

	if useTLS {
		server.AutoTLSManager.Cache = autocert.DirCache("/var/www/.cache")
	}

	server.Pre(middleware.NonWWWRedirect())
	server.Use(RemoveExtraSlashesMiddleware())
	server.Use(middleware.Logger())
	server.Use(middleware.BodyLimit(bodyLimit))

	if behindProxy {
		server.Use(echo.WrapMiddleware(handlers.ProxyHeaders))
	} else {
		server.Use(middleware.Gzip())
		server.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
			Skipper: func(context echo.Context) bool {
				sessionCookie, err := context.Cookie("session")
				if sessionCookie != nil && err == nil {
					return false
				}
				return true
			},
			TokenLength:  32,
			TokenLookup:  "header:" + echo.HeaderXCSRFToken,
			ContextKey:   "csrf",
			CookieName:   "csrf",
			CookieMaxAge: 86400,
			CookiePath:   "/",
			CookieSecure: useTLS,
		}))
	}

	if useTLS {
		server.Use(middleware.HTTPSRedirect())
		server.Use(middleware.SecureWithConfig(middleware.SecureConfig{
			XFrameOptions:         "DENY",
			HSTSMaxAge:            31536000,
			ContentSecurityPolicy: "default-src https:",
		}))
	} else {
		server.Use(middleware.SecureWithConfig(middleware.SecureConfig{
			XFrameOptions: "DENY",
			HSTSMaxAge:    31536000,
		}))
	}

	return &server
}
