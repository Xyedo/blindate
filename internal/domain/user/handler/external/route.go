package external

import (
	"github.com/labstack/echo/v4"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/infrastructure/auth"
)

func (u *User) Route(e *echo.Group) {
	users := e.Group("/users")

	users.GET("/:id/detail", u.getUserDetailByIdHandler)
	users.PUT("/:id/detail/photos", u.putUserDetailPhotoHandler, u.matchRequestParamId)
	users.POST("/:id/detail", u.postUserDetailHandler, u.matchRequestParamId)
	users.PATCH("/:id/detail", u.patchUserDetailByIdHandler, u.matchRequestParamId)

	users.POST("/:id/detail/interest", u.postInterestHandler, u.matchRequestParamId)
	users.PATCH("/:id/detail/interest", u.patchInterestHandler, u.matchRequestParamId)
	users.POST("/:id/detail/interest/delete", u.postDeleteInterestHandler, u.matchRequestParamId)
}

func (u *User) matchRequestParamId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		paramUserId := c.Param("id")

		requestId := c.Request().Context().Value(auth.RequestId).(string)
		if requestId != paramUserId {
			return apperror.NotFound(apperror.Payload{})
		}

		return next(c)
	}
}
