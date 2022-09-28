package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/internal/tokenizer"
)

func validateUser(jwtService tokenizer.Jwt) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		accessToken := strings.TrimPrefix(authorization, "Bearer ")

		id, err := jwtService.ValidateAccessToken(accessToken)
		if err != nil {
			if errors.Is(err, tokenizer.ErrTokenExpired) {
				errExpiredAccesToken(c)
				return
			}
			if errors.Is(err, tokenizer.ErrNotValidCredential) {
				errAccesTokenInvalid(c)
				return
			}
		}

		var url struct {
			Id string `uri:"id" binding:"required,uuid"`
		}
		err = c.ShouldBindUri(&url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "must have uuid in uri!",
			})
			return
		}

		if url.Id != id {
			errorInvalidIdTokenResponse(c)
			return
		}
		c.Set("userId", url.Id)
		c.Next()
	}
}
