package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/domain/match/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func CreateCandidateMatchsById(ctx context.Context, conn pg.Querier, userId string, candidateMatchsIds []string) error {
	rowAffected, err := conn.CopyFrom(ctx,
		pgx.Identifier{"match"},
		[]string{
			"request_from",
			"request_to",
			"request_status",
			"created_at",
			"updated_at",
			"version",
		},
		pgx.CopyFromSlice(len(candidateMatchsIds), func(i int) ([]any, error) {
			return []any{
				userId,
				candidateMatchsIds[i],
				entities.MatchStatusUnknown,
				time.Now(),
				time.Now(),
				1,
			}, nil
		}),
	)

	if err != nil {
		return err
	}

	if rowAffected != int64(len(candidateMatchsIds)) {
		return errors.New("something went wrong")
	}

	return nil

}

func FindUserMatchByStatus(ctx context.Context, userId string, status entities.MatchStatus, limit, page int) ([]entities.MatchUser, error) {
	const findUserMatchByStatus = `
	SELECT 
	FROM match m
	JOIN account_detail ad ON
		ad.account_id = m.request_from OR
		ad.account_id = m.request_to
	WHERE 
		ad.account_id = $1 AND
		m.status = $2
	`

	return nil, nil
}
