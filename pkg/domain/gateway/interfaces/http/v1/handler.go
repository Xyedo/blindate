package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/internal/security"
	httpMiddleware "github.com/xyedo/blindate/pkg/infrastructure/http/middleware"
)

func (h *gatewayH) Handler(globalRoutes *gin.RouterGroup, jwt *security.Jwt) {
	globalRoutes.GET("/ws", httpMiddleware.AuthToken(jwt), h.wsHandler)
}
