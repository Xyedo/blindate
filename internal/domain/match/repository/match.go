package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/common/ids"
	matchEntities "github.com/xyedo/blindate/internal/domain/match/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
	"github.com/xyedo/blindate/pkg/optional"
)

func CreateCandidateMatchsById(ctx context.Context, conn pg.Querier, userId string, candidateUserIds []string) error {
	if len(candidateUserIds) == 0 {
		return apperror.NotFound(apperror.Payload{
			Status:  matchEntities.ErrCodeMatchCandidateEmpty,
			Message: "empty candidate, pls try again later!",
		})
	}

	rowAffected, err := conn.CopyFrom(ctx,
		pgx.Identifier{"match"},
		[]string{
			"id",
			"request_from",
			"request_to",
			"request_status",
			"created_at",
			"updated_at",
			"updated_by",
			"version",
		},
		pgx.CopyFromSlice(len(candidateUserIds), func(i int) ([]any, error) {
			return []any{
				ids.Match(),
				userId,
				candidateUserIds[i],
				matchEntities.MatchStatusUnknown,
				time.Now(),
				time.Now(),
				nil,
				1,
			}, nil
		}),
	)

	if err != nil {
		return err
	}

	if rowAffected != int64(len(candidateUserIds)) {
		return errors.New("something went wrong")
	}

	return nil

}

func FindMatchsByStatus(ctx context.Context, conn pg.Querier, payload matchEntities.FindUserMatchByStatus) (matchEntities.Matchs, error) {
	const findUserMatchByStatus = `
	SELECT 
		m.id,
		m.request_from,
		m.request_to,
		m.request_status,
		m.accepted_at,
		m.reveal_status,
		m.revealed_declined_count,
		m.revealed_at,
		m.created_at,
		m.updated_at,
		m.updated_by,
		m.version
	FROM match m
	JOIN account_detail ad ON
	ad.account_id = m.request_from OR
	ad.account_id = m.request_to
	WHERE 
	ad.account_id = ?  AND
	m.request_status IN (?) AND
	CASE 
		WHEN m.request_status = 'REQUESTED' THEN m.updated_by != ? 
		ELSE TRUE 
	END
	LIMIT ?
	OFFSET ?
	`

	offset := payload.Limit*payload.Page - payload.Limit
	query, args, err := pg.In(
		findUserMatchByStatus,
		payload.UserId, payload.Statuses, payload.UserId, payload.Limit, offset,
	)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matchs := make([]matchEntities.Match, 0)
	for rows.Next() {
		var match matchEntities.Match
		var revealStatus optional.String
		err := rows.Scan(
			&match.Id,
			&match.RequestFrom,
			&match.RequestTo,
			&match.RequestStatus,
			&match.AcceptedAt,
			&revealStatus,
			&match.RevealedDeclinedCount,
			&match.RevealedAt,
			&match.CreatedAt,
			&match.UpdatedAt,
			&match.UpdatedBy,
			&match.Version,
		)
		if err != nil {
			return nil, err
		}
		revealStatus.If(func(s string) {
			match.RevealStatus.Set(matchEntities.MatchStatus(s))
		})

		matchs = append(matchs, match)
	}

	return matchs, nil
}

func GetMatchById(ctx context.Context, conn pg.Querier, id string, opts ...matchEntities.GetMatchOption) (matchEntities.Match, error) {
	const getMatchById = `
	SELECT 
		id,
		request_from,
		request_to,
		request_status,
		accepted_at,
		reveal_status,
		revealed_declined_count,
		revealed_at,
		created_at,
		updated_at,
		updated_by,
		version
	FROM match
	WHERE id = $1 
	`

	query := getMatchById
	if len(opts) > 0 && opts[0].PessimisticLocking {
		query += "\nFOR UPDATE"
	}

	var match matchEntities.Match
	var revealStatus optional.String
	err := conn.
		QueryRow(ctx, query, id).
		Scan(
			&match.Id,
			&match.RequestFrom,
			&match.RequestTo,
			&match.RequestStatus,
			&match.AcceptedAt,
			&revealStatus,
			&match.RevealedDeclinedCount,
			&match.RevealedAt,
			&match.CreatedAt,
			&match.UpdatedAt,
			&match.UpdatedBy,
			&match.Version,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return matchEntities.Match{}, apperror.NotFound(apperror.Payload{
				Error:  err,
				Status: matchEntities.ErrCodeMatchNotFound,
			})
		}

		return matchEntities.Match{}, err
	}

	revealStatus.If(func(s string) {
		match.RevealStatus.Set(matchEntities.MatchStatus(s))
	})

	return match, nil
}

func UpdateMatch(ctx context.Context, conn pg.Querier, match matchEntities.Match) error {
	const updateMatch = `
	UPDATE match SET
		request_from = $2,
		request_to = $3,
		request_status = $4,
		accepted_at = $5,
		reveal_status = $6,
		revealed_declined_count = $7,
		revealed_at =$8,
		created_at =$9,
		updated_at = $10,
		updated_by = $11,
		version = $12
	WHERE id = $1
	RETURNING id 
	`

	var revealStatus optional.String
	match.RevealStatus.If(func(ms matchEntities.MatchStatus) {
		revealStatus.Set(string(ms))
	})
	var returnedId string
	err := conn.
		QueryRow(ctx,
			updateMatch,
			match.Id,
			match.RequestFrom,
			match.RequestTo,
			string(match.RequestStatus),
			match.AcceptedAt,
			revealStatus,
			match.RevealedDeclinedCount,
			match.RevealedAt,
			match.CreatedAt,
			match.UpdatedAt,
			match.UpdatedBy,
			match.Version,
		).Scan(&returnedId)
	if err != nil {
		return err
	}

	if returnedId != match.Id {
		return errors.New("invalid")
	}

	return nil

}
