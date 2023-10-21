package handler

import "github.com/labstack/echo/v4"

func Route(e *echo.Echo) {
	user := e.Group("user")
	user.POST("/webhook", handleClrekWebhook)
}
