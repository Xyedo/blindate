package internal

import "github.com/labstack/echo/v4"

func Route(e *echo.Group) {
	user := e.Group("/users")
	user.POST("/event", handleEventWebhook)
}
