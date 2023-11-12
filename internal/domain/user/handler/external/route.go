package external

import (
	"github.com/labstack/echo/v4"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/infrastructure/auth"
)

func Route(e *echo.Group) {
	users := e.Group("/users")

	users.GET("/:id/detail", getUserDetailByIdHandler)
	users.PUT("/:id/detail/photo", putUserDetailPhotoHandler, matchRequestParamId)
	users.POST("/:id/detail", postUserDetailHandler, matchRequestParamId)
	users.PATCH("/:id/detail", patchUserDetailByIdHandler, matchRequestParamId)

	users.POST("/:id/detail/interest", postInterestHandler, matchRequestParamId)
	users.PATCH("/:id/detail/interest", patchInterestHandler, matchRequestParamId)
	users.POST("/:id/detail/interest/delete", postDeleteInterestHandler, matchRequestParamId)

}

func matchRequestParamId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		paramUserId := c.Param("id")

		requestId := c.Request().Context().Value(auth.RequestId).(string)

		if requestId != paramUserId {
			return apperror.NotFound(apperror.Payload{})
		}
		return next(c)
	}
}
