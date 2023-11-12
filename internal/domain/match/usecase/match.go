package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/domain/match/repository"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	userrepository "github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func CreateCandidateMatch(ctx context.Context, requestId string) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		user, err := userrepository.GetUserDetailById(ctx, tx, requestId)
		if err != nil {
			return err
		}

		closestUserIds, err := userrepository.FindNonMatchClosestUser(ctx, tx, entities.FindClosestUser{
			UserId: user.UserId,
			Geog:   user.Geog,
			Page:   1,
			Limit:  20,
		})
		hashset := make(map[string]struct{})
		hashset["asd"] = struct{}{}

		if err != nil {
			return err
		}

		return repository.CreateCandidateMatchsById(ctx, tx, requestId, closestUserIds)
	})
}
