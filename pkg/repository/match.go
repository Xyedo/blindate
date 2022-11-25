package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/entity"
)

type Match interface {
	InsertNewMatch(fromUserId, toUserId string) (string, error)
	SelectMatchByUserId(userId string) ([]entity.Match, error)
	UpdateMatchById(matchEntity entity.Match) error
	GetMatchById(matchId string) (*entity.Match, error)
}

func NewMatch(conn *sqlx.DB) *match {
	return &match{
		conn: conn,
	}
}

type match struct {
	conn *sqlx.DB
}

func (m *match) InsertNewMatch(fromUserId, toUserId string) (string, error) {
	query := `
	INSERT INTO match(
		request_from, 
		request_to, 
		created_at
		)
	VALUES($1,$2,$3)
	RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var matchId string
	err := m.conn.GetContext(ctx, &matchId, query, fromUserId, toUserId, time.Now())
	if err != nil {
		return "", err
	}
	return matchId, nil
}

func (m *match) SelectMatchByUserId(userId string) ([]entity.Match, error) {
	query := `
	SELECT 
		id,
		request_from, 
		request_to, 
		request_status, 
		created_at, 
		accepted_at, 
		reveal_status, 
		revealed_at
	FROM match
	WHERE request_from = $1 
		OR request_to = $1
	ORDER BY created_at ASC
	LIMIT 20`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	matchs := make([]entity.Match, 0)
	err := m.conn.SelectContext(ctx, &matchs, query, userId)
	if err != nil {
		return nil, err
	}
	return matchs, nil

}

func (m *match) GetMatchById(matchId string) (*entity.Match, error) {
	query := `
		SELECT
			id,
			request_from,
			request_to,
			request_status,
			created_at,
			accepted_at,
			reveal_status,
			revealed_at
		FROM match
		WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var matchEntity entity.Match
	err := m.conn.GetContext(ctx, &matchEntity, query, matchId)
	if err != nil {
		return nil, err
	}
	return &matchEntity, err
}
func (m *match) UpdateMatchById(matchEntity entity.Match) error {
	query := `
	UPDATE match SET
		request_status=$1, 
		accepted_at=$2, 
		reveal_status=$3, 
		revealed_at=$4
	WHERE id = $5
	RETURNING id`
	args := []any{
		matchEntity.RequestStatus,
		matchEntity.AcceptedAt,
		matchEntity.RevealStatus,
		matchEntity.RevealedAt,
		matchEntity.Id,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.conn.GetContext(ctx, &matchEntity.Id, query, args...)
	if err != nil {
		return err
	}
	return nil
}
