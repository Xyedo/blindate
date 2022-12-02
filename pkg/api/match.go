package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

type matchSvc interface {
	FindUserToMatch(userId string) ([]domain.BigUser, error)
	PostNewMatch(fromUserId, toUserId string, matchStatus domain.MatchStatus) (string, error)
	GetMatchReqToUserId(userId string) ([]domain.MatchUser, error)
	RequestChange(matchId string, matchStatus domain.MatchStatus) error
	RevealChange(matchId string, matchStatus domain.MatchStatus) error
}

func NewMatch(matchSvc matchSvc) *match {
	return &match{
		matchSvc: matchSvc,
	}
}

type match struct {
	matchSvc matchSvc
}

func (m *match) getNewMatchCandidateHandler(c *gin.Context) {
	userId := c.GetString("userId")
	res, err := m.matchSvc.FindUserToMatch(userId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTooLongAccessingDB):
			errResourceConflictResp(c)
		case errors.Is(err, domain.ErrResourceNotFound):
			errNotFoundResp(c, "userId is not found")
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"newCandidateMatchs": res,
		},
	})

}
func (m *match) postNewMatchHandler(c *gin.Context) {
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
	userId := c.GetString("userId")
	matchId, err := m.matchSvc.PostNewMatch(userId, input.ToUserId, domain.MatchStatus(input.MatchStatus))
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTooLongAccessingDB):
			errResourceConflictResp(c)
		case errors.Is(err, domain.ErrRefNotFound23503):
			errNotFoundResp(c, "your toUserId is not found in our db")
		case errors.Is(err, domain.ErrUniqueConstraint23505):
			errUnprocessableEntityResp(c, "already matched")
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"matchId": matchId,
		},
	})

}
func (m *match) getAllMatchHandler(c *gin.Context) {
	userId := c.GetString("userId")
	matcheds, err := m.matchSvc.GetMatchReqToUserId(userId)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccessingDB) {
			errResourceConflictResp(c)
			return
		}
		errServerResp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"matchs": matcheds,
		},
	})
}

func (m *match) putRequestHandler(c *gin.Context) {
	var input struct {
		Request string `form:"request" json:"request" binding:"required,oneof=requested declined accepted"`
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
	matchId := c.GetString("matchId")
	err := m.matchSvc.RequestChange(matchId, domain.MatchStatus(input.Request))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidMatchStatus):
			errUnprocessableEntityResp(c, "invalid matchStatus, please refer to the documentation")
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
	matchId := c.GetString("matchId")

	err := m.matchSvc.RevealChange(matchId, domain.MatchStatus(input.Reveal))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidMatchStatus):
			errUnprocessableEntityResp(c, "invalid matchStatus, please refer to the documentation")
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
