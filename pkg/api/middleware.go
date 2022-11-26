package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

const (
	authorizationHeaderKey = "Authorization"
)

func validateUser(jwtService jwtSvc) gin.HandlerFunc {
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
				errUnauthorizedResp(c, "token is expired, please login!")
			}
			if errors.Is(err, domain.ErrNotMatchCredential) {
				errUnauthorizedResp(c, "token is invalid, please login!")
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
			errForbiddenResp(c, "you should not access this resoures")
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
				"message": "required,must have uuid in uri!",
			})
			return
		}
		c.Set("interestId", url.InterestId)
		c.Next()
	}
}
func validateConversation() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url struct {
			ConversationId string `uri:"conversationId" binding:"required,uuid"`
		}
		err := c.ShouldBindUri(&url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "required,must have uuid in uri!",
			})
			return
		}
		c.Set("convId", url.ConversationId)
		c.Next()
	}

}

func validateChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url struct {
			ChatId string `uri:"chatId" binding:"required,uuid"`
		}
		err := c.ShouldBindUri(&url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "required,must have uuid in uri!",
			})
			return
		}
		c.Set("chatId", url.ChatId)
		c.Next()
	}

}
func validateMatch() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url struct {
			MatchId string `uri:"matchId" binding:"required,uuid"`
		}
		err := c.ShouldBindUri(&url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "required,must have uuid in uri!",
			})
			return
		}
		c.Set("matchId", url.MatchId)
		c.Next()
	}

}
