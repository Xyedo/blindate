package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	matchEntities "github.com/xyedo/blindate/internal/domain/match/entities"
	"github.com/xyedo/blindate/internal/domain/match/repository"
	matchRepo "github.com/xyedo/blindate/internal/domain/match/repository"
	userEntities "github.com/xyedo/blindate/internal/domain/user/entities"
	userRepo "github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func CreateCandidateMatch(ctx context.Context, requestId string) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		user, err := userRepo.GetUserDetailById(ctx, tx, requestId)
		if err != nil {
			return err
		}

		closestUserIds, err := userRepo.FindNonMatchClosestUser(ctx, tx, userEntities.FindClosestUser{
			UserId: user.UserId,
			Geog:   user.Geog,
			Page:   1,
			Limit:  20,
		})

		if err != nil {
			return err
		}

		return repository.CreateCandidateMatchsById(ctx, tx, requestId, closestUserIds)
	})
}

func IndexMatchCandidate(ctx context.Context, requestId string, limit, page int) ([]matchEntities.MatchUser, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	user, err := userRepo.GetUserDetailById(ctx, conn, requestId)
	if err != nil {
		return nil, err
	}

	matchUser, err := matchRepo.FindUserMatchByStatus(ctx, conn,
		matchEntities.FindUserMatchByStatus{
			UserId: user.UserId,
			Status: matchEntities.MatchStatusUnknown,
			Limit:  limit,
			Page:   page,
		},
	)
	if err != nil {
		return nil, err
	}

	err = matchUser.CalculateDistance(user)
	if err != nil {
		return nil, err
	}

	return matchUser, nil
}
