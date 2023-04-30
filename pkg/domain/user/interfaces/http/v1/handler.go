package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/internal/security"
	httpMiddleware "github.com/xyedo/blindate/pkg/infrastructure/http/middleware"
)

func (u *userH) Handler(globalRouter *gin.RouterGroup, jwt *security.Jwt) {
	users := globalRouter.Group("/users")

	users.POST("/", u.postUserHandler)

	userWithAuth := users.Group("/:userId", httpMiddleware.AuthToken(jwt), httpMiddleware.ValidateUserId())
	{
		userWithAuth.GET("/", u.getUserByIdHandler)
		userWithAuth.PATCH("/", u.patchUserByIdHandler)
		userWithAuth.PUT("/profile-picture", u.putUserImageProfileHandler)
	}
}
