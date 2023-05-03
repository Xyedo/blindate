package httpMiddleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/internal/security"
	"github.com/xyedo/blindate/pkg/common/constant"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

const (
	authorizationHeaderKey = "Authorization"
)

func AuthToken(jwt *security.Jwt) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid Authorization Header format",
			})
			return
		}

		if !strings.EqualFold("Bearer", fields[0]) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": fmt.Sprintf("Unsupported Authorization type %s", fields[0]),
			})
			return
		}

		accessToken := fields[1]
		id, err := jwt.ValidateAccessToken(accessToken)
		if err != nil {
			httperror.HandleError(c, err)
			return
		}

		c.Set(constant.KeyRequestUserId, id)
	}
}
func ValidateUserId() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url struct {
			UserId string `uri:"userId" binding:"required,uuid4"`
		}
		err := c.ShouldBindUri(&url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "must have uuid in uri!",
			})
			return
		}
		userId := c.GetString(constant.KeyRequestUserId)
		if url.UserId != userId {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "you should not access this resource",
			})
			return
		}
		c.Next()
	}
}
