package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func errValidationResp(c *gin.Context, errMap map[string]string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"status":  "fail",
		"message": "please refer to the documentation",
		"errors":  errMap,
	})
}

func errServerResp(c *gin.Context, err error) {
	log.Println(err.Error())
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"status":  "fail",
		"message": "the server encountered a problem and could not process your request",
	})
}
func errResourceConflictResp(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusConflict, gin.H{
		"status":  "fail",
		"message": "request conflicted, please try again",
	})
}
func errNotFoundResp(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"status":  "fail",
		"message": message,
	})
}

func errForbiddenResp(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"status":  "fail",
		"message": message,
	})
}
func errBadRequestResp(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status":  "fail",
		"message": message,
	})
}

func errUnprocessableEntityResp(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
		"status":  "fail",
		"message": message,
	})
}
