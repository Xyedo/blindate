package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func errorJSONBindingResponse(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status":  "fail",
		"message": err.Error(),
	})
}
func errorValidationResponse(c *gin.Context, errMap map[string]string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"status":  "fail",
		"message": "please refer to the documentation",
		"errors":  errMap,
	})
}

func errorServerResponse(c *gin.Context, err error) {
	log.Println(err.Error())
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"status":  "fail",
		"message": "the server encountered a problem and could not process your request",
	})
}

func errorInvalidCredsResponse(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  "fail",
		"message": message,
	})
}
func errorDeadLockResponse(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusConflict, gin.H{
		"status":  "fail",
		"message": "request conflicted, please try again",
	})
}
func errorResourceNotFound(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"status":  "fail",
		"message": message,
	})
}

func errCookieNotFound(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"status":  "fail",
		"message": "Cookie not found in your browser, must be login",
	})
}

func errExpiredAccesToken(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  "fail",
		"message": "token is expired, please login!",
	})
}
func errAccesTokenInvalid(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  "fail",
		"message": "token is invalid, please login!",
	})
}
func errorInvalidIdTokenResponse(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"status":  "fail",
		"message": "you should not access this resoures",
	})
}
