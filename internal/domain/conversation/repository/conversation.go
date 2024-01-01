package repository

import (
	"context"
	"errors"

	"github.com/xyedo/blindate/internal/domain/conversation/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func CreateConversation(ctx context.Context, conn pg.Querier, payload entities.Conversation) error {
	const createConversation = `
		INSERT INTO conversations (
			match_id,
			chat_rows,
			day_pass,
			created_at,
			updated_at,
			version
		)
		VALUES($1,$2,$3,$4,$5,$6)
		returning match_id
	`
	var returningMatchId string
	err := conn.
		QueryRow(ctx, createConversation,
			payload.MatchId,
			payload.ChatRows,
			payload.DayPass,
			payload.CreatedAt,
			payload.UpdatedAt,
			payload.Version,
		).Scan(&returningMatchId)
	if err != nil {
		return err
	}

	if returningMatchId != payload.MatchId {
		return errors.New("som wen wong")
	}

	return nil
}
