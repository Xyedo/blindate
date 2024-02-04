package dtos

import (
	"github.com/xyedo/blindate/internal/domain/match/entities"
	userdtos "github.com/xyedo/blindate/internal/domain/user/dtos"
)

type ShowMatchResponse struct {
	Id       string  `json:"id"`
	Status   string  `json:"status"`
	Distance float64 `json:"distance"`
	userdtos.UserDetail
}

func NewShowMatchResponse(matchUser entities.MatchUser) ShowMatchResponse {
	var status string
	switch matchUser.Status {
	case entities.MatchStatusRequested:
		status = "liked"
	case entities.MatchStatusAccepted:
		status = "accepted"
	default:
		panic("invalid status")
	}
	return ShowMatchResponse{
		Id:         matchUser.MatchId,
		Status:     status,
		Distance:   matchUser.Distance,
		UserDetail: userdtos.NewUserDetailResponse(matchUser.UserDetail),
	}
}
