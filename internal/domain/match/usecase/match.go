package usecase

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	matchEntities "github.com/xyedo/blindate/internal/domain/match/entities"
	"github.com/xyedo/blindate/internal/domain/match/repository"
	matchRepo "github.com/xyedo/blindate/internal/domain/match/repository"
	userEntities "github.com/xyedo/blindate/internal/domain/user/entities"
	userRepo "github.com/xyedo/blindate/internal/domain/user/repository"
	userUsecase "github.com/xyedo/blindate/internal/domain/user/usecase"
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

	matchUserIds, err := matchRepo.FindMatchUserIdsByStatus(ctx, conn,
		matchEntities.FindUserMatchByStatus{
			UserId:   user.UserId,
			Statuses: []matchEntities.MatchStatus{matchEntities.MatchStatusUnknown, matchEntities.MatchStatusRequested},
			Limit:    limit,
			Page:     page,
		},
	)
	if err != nil {
		return nil, err
	}

	userDetails, err := userUsecase.GetUserDetails(ctx, matchUserIds)
	if err != nil {
		return nil, err
	}

	return matchEntities.NewMatchUsers(user.Geog, userDetails), nil
}

func TransitionRequestStatus(ctx context.Context, requestId, matchId string, swipe bool) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		requester, err := userRepo.GetUserDetailById(ctx, tx, requestId)
		if err != nil {
			return err
		}

		match, err := matchRepo.GetMatchById(ctx, tx, matchId, matchEntities.GetMatchOption{
			PessimisticLocking: true,
		})
		if err != nil {
			return err
		}

		err = match.ValidateResource(requester)
		if err != nil {
			return err
		}

		switch match.RequestStatus {
		case matchEntities.MatchStatusUnknown:
			if swipe {
				match.RequestStatus = matchEntities.MatchStatusRequested
			} else {
				match.RequestStatus = matchEntities.MatchStatusDeclined
			}

		case matchEntities.MatchStatusRequested:
			if match.UpdatedBy.MustGet() == requester.UserId {
				return apperror.BadPayload(apperror.Payload{
					Status:  matchEntities.ErrCodeMatchStatusInvalid,
					Message: "invalid status",
				})
			}

			if swipe {
				match.RequestStatus = matchEntities.MatchStatusAccepted
			} else {
				match.RequestStatus = matchEntities.MatchStatusDeclined
			}

		case matchEntities.MatchStatusDeclined, matchEntities.MatchStatusAccepted:
			return nil
		}

		match.UpdatedAt = time.Now()
		match.UpdatedBy.Set(requester.UserId)
		match.Version++

		return matchRepo.UpdateMatch(ctx, tx, match)
	})
}
