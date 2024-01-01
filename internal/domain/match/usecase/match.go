package usecase

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	conversationEntities "github.com/xyedo/blindate/internal/domain/conversation/entities"
	conversationRepo "github.com/xyedo/blindate/internal/domain/conversation/repository"
	"github.com/xyedo/blindate/internal/domain/match/entities"
	"github.com/xyedo/blindate/internal/domain/match/repository"
	matchRepo "github.com/xyedo/blindate/internal/domain/match/repository"
	userEntities "github.com/xyedo/blindate/internal/domain/user/entities"
	userRepo "github.com/xyedo/blindate/internal/domain/user/repository"
	userUsecase "github.com/xyedo/blindate/internal/domain/user/usecase"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
	"github.com/xyedo/blindate/pkg/pagination"
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
			Pagination: pagination.Pagination{
				Page:  1,
				Limit: 10,
			},
		})

		if err != nil {
			return err
		}

		return repository.CreateCandidateMatchsById(ctx, tx, requestId, closestUserIds)
	})
}

func IndexMatchCandidate(ctx context.Context, requestId string, limit, page int) ([]entities.MatchUser, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	requestUser, err := userRepo.GetUserDetailById(ctx, conn, requestId)
	if err != nil {
		return nil, err
	}

	matchs, err := matchRepo.FindMatchsByStatus(ctx, conn,
		entities.FindUserMatchByStatus{
			UserId: requestUser.UserId,
			Statuses: []entities.MatchStatus{
				entities.MatchStatusUnknown,
				entities.MatchStatusRequested,
			},
			Limit: limit,
			Page:  page,
		},
	)
	if err != nil {
		return nil, err
	}

	matchUserIds, matchUserIdToMatchId := matchs.ToUserIds(requestId)
	userDetails, err := userUsecase.GetUserDetails(ctx, matchUserIds)
	if err != nil {
		return nil, err
	}

	return entities.NewMatchUsers(
		requestUser,
		userDetails,
		matchUserIdToMatchId,
	), nil
}

func TransitionRequestStatus(ctx context.Context, requestId, matchId string, swipe bool) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		requester, err := userRepo.GetUserDetailById(ctx, tx, requestId)
		if err != nil {
			return err
		}

		match, err := matchRepo.GetMatchById(ctx, tx, matchId, entities.GetMatchOption{
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
		case entities.MatchStatusUnknown:
			if swipe {
				match.RequestStatus = entities.MatchStatusRequested
			} else {
				match.RequestStatus = entities.MatchStatusDeclined
			}

		case entities.MatchStatusRequested:
			if match.UpdatedBy.MustGet() == requester.UserId {
				return apperror.BadPayload(apperror.Payload{
					Status:  entities.ErrCodeMatchStatusInvalid,
					Message: "invalid status",
				})
			}

			if swipe {
				match.RequestStatus = entities.MatchStatusAccepted

				err = conversationRepo.CreateConversation(ctx, tx, conversationEntities.Conversation{
					MatchId:   match.Id,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Version:   1,
				})
				if err != nil {
					return err
				}

			} else {
				match.RequestStatus = entities.MatchStatusDeclined
			}

		case entities.MatchStatusDeclined, entities.MatchStatusAccepted:
			return nil
		}

		match.UpdatedAt = time.Now()
		match.UpdatedBy.Set(requester.UserId)
		match.Version++

		return matchRepo.UpdateMatch(ctx, tx, match)
	})
}
