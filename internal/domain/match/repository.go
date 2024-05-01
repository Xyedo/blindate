package match

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/match/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

type Repository interface {
	CreateCandidateMatchsById(ctx context.Context, conn pg.Querier, userId string, candidateUserIds []string) error
	FindMatchsByStatus(ctx context.Context, conn pg.Querier, payload entities.FindUserMatchByStatus) (entities.Matchs, bool, error)
	GetMatchById(ctx context.Context, conn pg.Querier, id string, opts ...entities.GetMatchOption) (entities.Match, error)
	UpdateMatch(ctx context.Context, conn pg.Querier, match entities.Match) error
}
