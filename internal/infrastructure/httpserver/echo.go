package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	userhandler "github.com/xyedo/blindate/internal/domain/user/handler"
	"github.com/xyedo/blindate/internal/infrastructure"
)

func NewEcho() *Server {
	e := echo.New()
	e.HTTPErrorHandler = EchoErrorHandler

	e.Use(middleware.Recover(), middleware.CORS(), middleware.BodyLimit("4M"), middleware.ContextTimeout(3*time.Second))

	e.GET("/healthcheck", func(c echo.Context) error {
		return nil
	})

	routeHandler(e)

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", infrastructure.Config.Host, infrastructure.Config.Port),
			Handler: e,
		},
	}
}
func routeHandler(e *echo.Echo) {
	userhandler.Route(e)
}
