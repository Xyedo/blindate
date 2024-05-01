package usecase

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	geo "github.com/paulmach/go.geo"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	conversationEntities "github.com/xyedo/blindate/internal/domain/conversation/entities"
	"github.com/xyedo/blindate/internal/domain/match"
	"github.com/xyedo/blindate/internal/domain/match/entities"
	"github.com/xyedo/blindate/internal/domain/match/repository"
	userEntities "github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
	"github.com/xyedo/blindate/pkg/pagination"
)

func New(repo repository.Match, userUsecase match.UserUsecase) *Match {
	return &Match{
		repo:        repo,
		userUsecase: userUsecase,
	}
}

type Match struct {
	repo                repository.Match
	userUsecase         match.UserUsecase
	conversationUsecase match.MatchUsecase
}

var _ match.Usecase = &Match{}

func (uc *Match) CreateCandidateMatch(ctx context.Context, requestId string) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		user, err := uc.userUsecase.GetUserDetailByUserId(ctx, tx, requestId)
		if err != nil {
			return err
		}

		closestUserIds, err := uc.userUsecase.FindNonMatchClosestUserIds(ctx, tx, userEntities.FindClosestUser{
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

		return uc.repo.CreateCandidateMatchsById(ctx, tx, requestId, closestUserIds)
	})
}

func (uc *Match) IndexMatch(ctx context.Context, requestId string, payload entities.IndexMatch) ([]entities.MatchUser, bool, error) {
	var (
		matchUsers []entities.MatchUser
		hasNext    bool
	)
	err := pg.TransactionWithRetry(ctx,
		pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly},
		func(tx pg.Querier) error {
			requestUser, err := uc.userUsecase.GetUserDetailByUserId(ctx, tx, requestId)
			if err != nil {
				return err
			}

			matchs, h, err := uc.repo.FindMatchsByStatus(ctx, tx,
				entities.FindUserMatchByStatus{
					UserId:     requestUser.UserId,
					Statuses:   payload.MatchStatuses(),
					Pagination: payload.Pagination,
				},
			)
			if err != nil {
				return err
			}
			hasNext = h

			matchUserIds, matchUserIdToMatchId := matchs.ToUserIds(requestId)
			userDetails, err := uc.userUsecase.GetUserDetails(ctx, tx, matchUserIds)
			if err != nil {
				return err
			}

			matchUsers = entities.NewMatchUsers(
				requestUser,
				userDetails,
				matchUserIdToMatchId,
			)

			return nil
		},
	)

	if err != nil {
		return nil, false, err
	}
	return matchUsers, hasNext, nil

}

func (uc *Match) GetMatchById(ctx context.Context, requestId, matchId string) (entities.MatchUser, error) {
	var (
		requestUser, recepientUser userEntities.UserDetail
		match                      entities.Match
	)
	err := pg.TransactionWithRetry(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly}, func(tx pg.Querier) error {
		user, err := uc.userUsecase.GetUserDetailByUserId(ctx, tx, requestId)
		if err != nil {
			return err
		}
		requestUser = user

		returnedMatch, err := uc.repo.GetMatchById(ctx, tx, matchId)
		if err != nil {
			return err
		}
		match = returnedMatch

		recepientId, err := match.ValidateResource(requestId)
		if err != nil {
			return err
		}

		err = match.ValidateShow(requestId)
		if err != nil {
			return err
		}

		user, err = uc.userUsecase.GetUserDetailByUserId(ctx, tx, recepientId)
		if err != nil {
			return err
		}
		recepientUser = user

		return nil
	})
	if err != nil {
		return entities.MatchUser{}, err
	}

	return entities.MatchUser{
		MatchId: matchId,
		Status:  match.RequestStatus,
		Distance: geo.
			NewPointFromLatLng(recepientUser.Geog.Lat, recepientUser.Geog.Lng).
			DistanceFrom(
				geo.NewPoint(requestUser.Geog.Lat, requestUser.Geog.Lng),
			),
		UserDetail: recepientUser,
	}, nil

}

func (uc *Match) TransitionRequestStatus(ctx context.Context, requestId, matchId string, swipe bool) error {
	return pg.TransactionWithRetry(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, func(tx pg.Querier) error {
		requester, err := uc.userUsecase.GetUserDetailByUserId(ctx, tx, requestId)
		if err != nil {
			return err
		}

		match, err := uc.repo.GetMatchById(ctx, tx, matchId)
		if err != nil {
			return err
		}

		_, err = match.ValidateResource(requestId)
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

				err = uc.conversationUsecase.CreateConversation(ctx, tx, conversationEntities.Conversation{
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

		return uc.repo.UpdateMatch(ctx, tx, match)
	})

}
