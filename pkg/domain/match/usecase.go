package match

import (
	matchDTOs "github.com/xyedo/blindate/pkg/domain/match/dtos"
)

type Usecase interface {
	FindMatchCandidate(string) ([]matchDTOs.FullUserDetail, error)
	CreateMatchFromCandidate(fromUserId, toUserId string, status matchDTOs.Status) (string, error)
	GetUserRequestedMatch(string) ([]matchDTOs.FullUserDetail, error)
	RequestChangeStateByMatchId(string, status matchDTOs.Status) error
	
}
