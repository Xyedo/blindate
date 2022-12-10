package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "Authorization"
)

func authToken(jwtSvc jwtSvc) gin.HandlerFunc {
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
		id, err := jwtSvc.ValidateAccessToken(accessToken)
		if err != nil {
			jsonHandleError(c, err)
			return
		}
		c.Set(keyUserId, id)
	}
}
func validateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url struct {
			UserId string `uri:"userId" binding:"required,uuid"`
		}
		err := c.ShouldBindUri(&url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "must have uuid in uri!",
			})
			return
		}
		userId := c.GetString(keyUserId)
		if url.UserId != userId {
			errForbiddenResp(c, "you should not access this resoures")
			return
		}
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
		c.Set(keyInterestId, url.InterestId)
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
		c.Set(keyConvId, url.ConversationId)
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
		c.Set(keyChatId, url.ChatId)
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
		c.Set(keyMatchId, url.MatchId)
		c.Next()
	}

}
