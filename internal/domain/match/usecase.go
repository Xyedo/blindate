package match

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/match/entities"
)

type Usecase interface {
	CreateCandidateMatch(ctx context.Context, requestId string) error
	GetMatchById(ctx context.Context, requestId string, matchId string) (entities.MatchUser, error)
	IndexMatch(ctx context.Context, requestId string, payload entities.IndexMatch) ([]entities.MatchUser, bool, error)
	TransitionRequestStatus(ctx context.Context, requestId string, matchId string, swipe bool) error
}
