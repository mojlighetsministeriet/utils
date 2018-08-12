package server

import (
	"net/http"
	"sort"
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

type Route struct {
	Path   string `json:"path"`
	Method string `json:"method"`
	Name   string `json:"name"`
}

type Routes []Route

func (routes Routes) Sort() {
	sorter := &routesSorter{routes: routes}
	sort.Sort(sorter)
}

type routesSorter struct {
	routes []Route
}

func (sorter *routesSorter) Len() int {
	return len(sorter.routes)
}

func (sorter *routesSorter) Swap(i, j int) {
	sorter.routes[i], sorter.routes[j] = sorter.routes[j], sorter.routes[i]
}

func (sorter *routesSorter) Less(i, j int) bool {
	return sorter.routes[i].Path < sorter.routes[j].Path
}

func (server *Server) UseTLS() bool {
	return server.useTLS
}

func (server *Server) Listen(address string) {
	server.addHelpResourceIfMissing()

	if server.useTLS {
		server.Logger.Fatal(server.StartAutoTLS(address))
	} else {
		server.Logger.Fatal(server.Start(address))
	}
}

func (server *Server) addHelpResourceIfMissing() {
	var registeredRoutes Routes
	helpIsMissing := true

	for _, route := range server.Routes() {
		if !strings.HasSuffix(route.Path, "/*") {
			registeredRoute := Route{
				Path:   route.Path,
				Method: route.Method,
				Name:   route.Name,
			}
			registeredRoutes = append(registeredRoutes, registeredRoute)
		}

		if route.Method == "GET" && route.Path == "/help" {
			helpIsMissing = false
		}
	}

	if helpIsMissing {
		registeredRoutes.Sort()

		server.GET("/help", func(context echo.Context) error {
			return context.JSON(http.StatusOK, registeredRoutes)
		})
	}
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
