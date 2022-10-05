package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/service"
)

const (
	authorizationHeaderKey = "Authorization"
)

func validateUser(jwtService service.Jwt) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader(authorizationHeaderKey)
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "Invalid Authorization Header format",
			})
			return
		}
		if !strings.EqualFold("Bearer", fields[0]) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": fmt.Sprintf("Unsupported Authorization type %s", fields[0]),
			})
			return
		}

		accessToken := fields[1]
		id, err := jwtService.ValidateAccessToken(accessToken)
		if err != nil {
			if errors.Is(err, service.ErrTokenExpired) {
				errExpiredAccesToken(c)
			}
			if errors.Is(err, service.ErrTokenNotValid) {
				errAccesTokenInvalid(c)
			}
			return
		}

		var url struct {
			UserId string `uri:"userId" binding:"required,uuid"`
		}
		err = c.ShouldBindUri(&url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "must have uuid in uri!",
			})
			return
		}

		if url.UserId != id {
			errorInvalidIdTokenResponse(c)
			return
		}
		c.Set("userId", url.UserId)
		c.Next()
	}
}

func validateInterest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url struct {
			InterestId string `uri:"interestId" binding:"required,uuid"`
		}
		err := c.ShouldBindUri(&url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "must have uuid in uri!",
			})
			return
		}
		c.Set("interestId", url.InterestId)
		c.Next()
	}

}
