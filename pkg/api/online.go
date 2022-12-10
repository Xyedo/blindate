package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

type onlineSvc interface {
	CreateNewOnline(userId string) error
	PutOnline(userId string, online bool) error
	GetOnline(userId string) (domain.Online, error)
}

func NewOnline(onlineSvc onlineSvc) *Online {
	return &Online{
		onlineSvc: onlineSvc,
	}
}

type Online struct {
	onlineSvc onlineSvc
}

func (o *Online) postUserOnlineHandler(c *gin.Context) {
	userId := c.GetString(keyUserId)

	err := o.onlineSvc.CreateNewOnline(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "user-online created",
	})

}
func (o *Online) getUserOnlineHandler(c *gin.Context) {
	userId := c.GetString(keyUserId)
	onlineUser, err := o.onlineSvc.GetOnline(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"online": onlineUser,
		},
	})
}
func (o *Online) putuserOnlineHandler(c *gin.Context) {
	userId := c.GetString(keyUserId)
	var input struct {
		Online *bool `json:"online" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"online": "required and should be boolean",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = o.onlineSvc.PutOnline(userId, *input.Online)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user-online updated",
	})
}
