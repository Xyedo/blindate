package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/internal/security"
	httpMiddleware "github.com/xyedo/blindate/pkg/infrastructure/http/middleware"
)

func (h *onlineH) Handler(globalRoute *gin.RouterGroup, jwt *security.Jwt) {
	online := globalRoute.Group("/users/:userId/online", httpMiddleware.AuthToken(jwt))
	online.POST("/", httpMiddleware.ValidateUserId(), h.postOnlineHandler)
	online.GET("/", h.getOnlineHandler)
}
