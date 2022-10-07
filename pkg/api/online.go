package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

func NewOnline(onlineSvc service.Online) *online {
	return &online{
		onlineSvc: onlineSvc,
	}
}

type online struct {
	onlineSvc service.Online
}

func (o *online) postUserOnlineHandler(c *gin.Context) {
	userId := c.GetString("userId")

	err := o.onlineSvc.CreateNewOnline(userId)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "userId is not found!")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainUserId) {
			errorBadRequest(c, "users is already registered")
			return
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "user-online created",
	})

}
func (o *online) getUserOnlineHandler(c *gin.Context) {
	userId := c.GetString("userId")
	onlineUser, err := o.onlineSvc.GetOnline(userId)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "user_id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"online": onlineUser,
		},
	})
}
func (o *online) putuserOnlineHandler(c *gin.Context) {
	userId := c.GetString("userId")
	var input struct {
		Online bool `json:"online" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Online": "required and should be boolean",
		})
		if errjson != nil {
			errorServerResponse(c, err)
			return
		}
		return
	}
	err = o.onlineSvc.PutOnline(userId, input.Online)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "user_id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user-online updated",
	})
}
