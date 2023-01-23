package match

import (
	matchEntity "github.com/xyedo/blindate/pkg/domain/match/entities"
)

type Repository interface {
	InsertNewMatch(fromUserId, toUserId string, reqStatus matchEntity.Status) (string, error)
	SelectMatchReqToUserId(userId string) ([]matchEntity.FullUserDTO, error)
	UpdateMatchById(matchEntity matchEntity.MatchDAO) error
	GetMatchById(matchId string) (matchEntity.MatchDAO, error)
}
