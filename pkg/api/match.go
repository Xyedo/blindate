package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

// TODO: Create match_test.go
type matchSvc interface {
	FindUserToMatch(userId string) ([]domain.BigUser, error)
	PostNewMatch(fromUserId, toUserId string, matchStatus domain.MatchStatus) (string, error)
	GetMatchReqToUserId(userId string) ([]domain.MatchUser, error)
	RequestChange(matchId string, matchStatus domain.MatchStatus) error
	RevealChange(matchId string, matchStatus domain.MatchStatus) error
}

func NewMatch(matchSvc matchSvc) *Match {
	return &Match{
		matchSvc: matchSvc,
	}
}

type Match struct {
	matchSvc matchSvc
}

func (m *Match) getNewUserToMatchHandler(c *gin.Context) {
	userId := c.GetString(keyUserId)
	res, err := m.matchSvc.FindUserToMatch(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"newCandidateMatchs": res,
		},
	})

}
func (m *Match) postNewMatchHandler(c *gin.Context) {
	var input struct {
		ToUserId    string `json:"toUserId" form:"toUserId" binding:"required,uuid"`
		MatchStatus string `json:"matchStatus" form:"matchStatus" binding:"required,oneof=requested declined"`
	}
	err := c.ShouldBind(&input)
	if err != nil {
		if errjson := jsonBindingErrResp(err, c, map[string]string{
			"toUserId":    "must be provided and must be valid uuid",
			"matchStatus": "must be provided and the values is one of `requested` or `declined`",
		}); errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	userId := c.GetString(keyUserId)
	matchId, err := m.matchSvc.PostNewMatch(userId, input.ToUserId, domain.MatchStatus(input.MatchStatus))
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"matchId": matchId,
		},
	})

}
func (m *Match) getAllMatchRequestedHandler(c *gin.Context) {
	userId := c.GetString(keyUserId)
	matcheds, err := m.matchSvc.GetMatchReqToUserId(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"matchs": matcheds,
		},
	})
}

func (m *Match) putRequestHandler(c *gin.Context) {
	var input struct {
		Request string `form:"request" json:"request" binding:"required,oneof=declined accepted"`
	}
	if err := c.ShouldBind(&input); err != nil {
		if errjson := jsonBindingErrResp(err, c, map[string]string{
			"request": "required and the value must be `accepted` or `declined`",
		}); errjson != nil {
			errServerResp(c, err)
			return
		}
		return

	}
	matchId := c.GetString(keyMatchId)
	err := m.matchSvc.RequestChange(matchId, domain.MatchStatus(input.Request))
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "request has been accepted",
	})
}
func (m *Match) putRevealHandler(c *gin.Context) {
	var input struct {
		Reveal string `form:"reveal" json:"reveal" binding:"required,oneof=requested accepted declined"`
	}
	if err := c.ShouldBind(&input); err != nil {
		if errjson := jsonBindingErrResp(err, c, map[string]string{
			"Request": "required and the value must be `requested` or `accepted` or `declined`",
		}); errjson != nil {
			errServerResp(c, err)
			return
		}
		return

	}
	matchId := c.GetString(keyMatchId)

	err := m.matchSvc.RevealChange(matchId, domain.MatchStatus(input.Reveal))
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "request has been accepted",
	})
}
