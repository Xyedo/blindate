package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

type matchSvc interface {
	FindNewMatch(userId string) ([]domain.User, error)
	AcceptRequest(matchId string) error
	RequestReveal(matchId string) error
	AcceptReveal(matchId string) error
}

func NewMatch(matchSvc matchSvc, interestSvc interestSvc, basicInfoSvc basicInfoSvc) *match {
	return &match{
		matchSvc:     matchSvc,
		interestSvc:  interestSvc,
		basicInfoSvc: basicInfoSvc,
	}
}

type match struct {
	matchSvc     matchSvc
	interestSvc  interestSvc
	basicInfoSvc basicInfoSvc
}

func (m *match) getNewMatchHandler(c *gin.Context) {
	var result []struct {
		domain.User
		domain.Interest
		domain.BasicInfo
	}
	userId := c.GetString("userId")
	matchUsers, err := m.matchSvc.FindNewMatch(userId)
	if err != nil {
		//TODO: better error handling
		errServerResp(c, err)
		return
	}
	for _, matchUser := range matchUsers {
		basicInfo, err := m.basicInfoSvc.GetBasicInfoByUserId(matchUser.ID)
		if err != nil {
			//TODO: better error handling
			errServerResp(c, err)
			return
		}

		interest, err := m.interestSvc.GetInterest(matchUser.ID)
		if err != nil {
			//TODO: better error handling
			errServerResp(c, err)
			return
		}
		result = append(result, struct {
			domain.User
			domain.Interest
			domain.BasicInfo
		}{
			User:      matchUser,
			Interest:  *interest,
			BasicInfo: *basicInfo,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"results": result,
		},
	})

}

func (m *match) putRequestHandler(c *gin.Context) {
	var input struct {
		Request bool `form:"request" json:"request" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		if errjson := jsonBindingErrResp(err, c, map[string]string{
			"Request": "required and the value must be boolean true`",
		}); errjson != nil {
			errServerResp(c, err)
			return
		}
		return

	}
	matchId := c.GetString("matchId")
	err := m.matchSvc.AcceptReveal(matchId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrResourceNotFound):
			errNotFoundResp(c, "matchId is not found")
		case errors.Is(err, domain.ErrTooLongAccessingDB):
			errResourceConflictResp(c)
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "request has been accepted",
	})
}
func (m *match) putRevealHandler(c *gin.Context) {
	var input struct {
		Reveal string `form:"reveal" json:"reveal" binding:"required,oneof=requested accepted"`
	}
	if err := c.ShouldBind(&input); err != nil {
		if errjson := jsonBindingErrResp(err, c, map[string]string{
			"Request": "required and the value must be `requested` or `accepted`",
		}); errjson != nil {
			errServerResp(c, err)
			return
		}
		return

	}
	matchId := c.GetString("matchId")

	var err error
	if input.Reveal == string(domain.Requested) {
		err = m.matchSvc.RequestReveal(matchId)
	} else {
		err = m.matchSvc.AcceptReveal(matchId)
	}
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrResourceNotFound):
			errNotFoundResp(c, "matchId is not found")
		case errors.Is(err, domain.ErrTooLongAccessingDB):
			errResourceConflictResp(c)
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "request has been accepted",
	})
}
