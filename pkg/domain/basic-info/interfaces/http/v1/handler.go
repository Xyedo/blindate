package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/internal/security"
	httpMiddleware "github.com/xyedo/blindate/pkg/infrastructure/http/middleware"
)

func (b *basicInfoH) Handler(globalRoute *gin.RouterGroup, jwt *security.Jwt) {
	basicInfo := globalRoute.Group("/users/:userId/basic-info", httpMiddleware.AuthToken(jwt))
	basicInfo.POST("/", httpMiddleware.ValidateUserId(), b.postBasicInfoHandler)
	basicInfo.GET("/", b.getBasicInfoHandler)
	basicInfo.PATCH("/", httpMiddleware.ValidateUserId(), b.patchBasicInfoHandler)
}
