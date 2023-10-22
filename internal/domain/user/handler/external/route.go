package external

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xyedo/blindate/internal/infrastructure/auth"
)

func Route(e *echo.Group) {
	users := e.Group("users")
	users.GET("/:id/basic-info", getBasicInfoById)
	users.POST("/:id/basic-info", postBasicInfo, matchRequestParamId)
	users.PATCH("/:id/basic-info", patchBasicInfoById, matchRequestParamId)
}

func matchRequestParamId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		paramUserId := c.Param("id")

		requestId := c.Request().Context().Value(auth.RequestId).(string)

		if requestId != paramUserId {
			return c.NoContent(http.StatusMethodNotAllowed)
		}
		return next(c)
	}
}
