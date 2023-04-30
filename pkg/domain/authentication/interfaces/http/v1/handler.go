package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
)

func (h *authH) Handler(globalRoutes *gin.RouterGroup) {
	auth := globalRoutes.Group("/auth")

	auth.POST("/", h.postAuthHandler)
	
	{
		authWithRefreshToken := auth.Group("/", withRefreshToken())
		authWithRefreshToken.PUT("/", h.putAuthHandler)
		authWithRefreshToken.DELETE("/", h.deleteAuthHandler)
	}

}

func withRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshTokenCookie, err := c.Cookie("refreshToken")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Cookie not found in your browser, must be login",
			})
			return
		}
		c.Set(constant.KeyRefreshToken, refreshTokenCookie)
	}
}
