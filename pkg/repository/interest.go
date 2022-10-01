package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/entity"
)

func NewInterest(db *sqlx.DB) *interest {
	return &interest{
		db,
	}
}

type interest struct {
	*sqlx.DB
}

func (i *interest) InsertNewInterest(intr *entity.Interest) (int64, error) {
	query := `
	INSERT INTO interests (
		user_id, 
		hobbies,
		movies_series,
		traveling, 
		sport, 
		bio, 
		spotify_connect, 
		created_at, 
		updated_at)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8,$8)`
	args := []any{
		intr.UserId,
		pq.Array(intr.Hobbies),
		pq.Array(intr.MoviesSeries),
		pq.Array(intr.Traveling),
		pq.Array(intr.Sport),
		intr.Bio,
		intr.SpotifyConnect,
		time.Now(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := i.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return row, nil
}
