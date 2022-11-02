package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/util"
)

func jsonBindingErrResp(err error, c *gin.Context, errorMap map[string]string) error {
	err1 := util.ReadJSONDecoderErr(err)
	if err1 != nil {
		errBadRequestResp(c, err1.Error())
		return nil
	}
	errMap := util.ReadValidationErr(err, errorMap)
	if errMap != nil {
		errValidationResp(c, errMap)
		return nil
	}
	return err
}

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

func errUnauthorizedResp(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  "fail",
		"message": message,
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
