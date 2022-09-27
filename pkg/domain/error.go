package domain

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrDuplicateEmail     = errors.New("database: duplicate email")
	ErrNotMatchCredential = errors.New("database: not match credential")
	ErrTooLongAccesingDB  = errors.New("database: too long accessing DB")
	ErrDuplicateToken     = errors.New("database: duplicate token")
	ErrUserNotFound       = errors.New("database: user not found")
)

func ErrorJSONBindingResponse(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status":  "fail",
		"message": err.Error(),
	})
}
func ErrorValidationResponse(c *gin.Context, errMap map[string]string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"status":  "fail",
		"message": "please refer to the documentation",
		"error":   errMap,
	})
}

func ErrorServerResponse(c *gin.Context, err error) {
	log.Println(err.Error())
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"status":  "fail",
		"message": "the server encountered a problem and could not process your request",
	})
}

func ErrorInvalidCredsResponse(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  "fail",
		"message": message,
	})
}
func ErrorRequestTimeout(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
		"status":  "fail",
		"message": "request timeout, please refer to the documentation",
	})
}
func ErrorResourceNotFound(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status":  "fail",
		"message": message,
	})
}
