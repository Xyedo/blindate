package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xyedo/blindate/internal/infrastructure"
)

func InternalApiKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		apiKey := c.Request().Header.Get("ApiKey")
		if apiKey != infrastructure.Config.Clerk.ApiKey {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"message": "invalid request header key/value",
			})
		}

		return next(c)
	}
}
