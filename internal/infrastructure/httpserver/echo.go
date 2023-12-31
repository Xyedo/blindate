package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	userHandler "github.com/xyedo/blindate/internal/domain/user/handler/external"
	internalUserHandler "github.com/xyedo/blindate/internal/domain/user/handler/intrnl"

	"github.com/xyedo/blindate/internal/infrastructure"
	echomiddleware "github.com/xyedo/blindate/internal/infrastructure/httpserver/echo-middleware"
)

func NewEcho() *Server {
	e := echo.New()
	e.HTTPErrorHandler = EchoErrorHandler

	e.Use(
		middleware.Recover(),
		middleware.CORS(),
		middleware.BodyLimit("4M"),
		middleware.ContextTimeout(3*time.Second))

	e.GET("/healthcheck", func(c echo.Context) error {
		return nil
	})

	apiv1 := e.Group("/v1")
	internalRouteHandler(apiv1)

	apiv1.Use(echomiddleware.Guard)

	userHandler.Route(apiv1)

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", infrastructure.Config.Host, infrastructure.Config.Port),
			Handler: e,
		},
	}
}

func internalRouteHandler(e *echo.Group) {
	internal := e.Group("/internal")
	internalUserHandler.Route(internal)
}
