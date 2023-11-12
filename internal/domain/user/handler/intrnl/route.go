package internal

import "github.com/labstack/echo/v4"

func Route(e *echo.Group) {
	user := e.Group("/users")
	user.POST("/event", handleEventWebhook)
}

func guardInternal(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
