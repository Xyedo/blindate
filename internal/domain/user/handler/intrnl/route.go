package internal

import (
	"github.com/labstack/echo/v4"
	middleware "github.com/xyedo/blindate/internal/infrastructure/httpserver/echo-middleware"
)

func Route(e *echo.Group) {
	user := e.Group("/users")
	user.POST("/event", handleEventWebhook, middleware.InternalApiKey)
}
