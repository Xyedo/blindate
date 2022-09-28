package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/entity"
)

func NewLocation(db *sqlx.DB) *location {
	return &location{
		db,
	}
}

type location struct {
	*sqlx.DB
}

func (l *location) InsertNewLocation(location *entity.Location) (int64, error) {
	query := `
		INSERT INTO locations(user_id, geog)
		VALUES($1, ST_GeomFromText($2))`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := l.ExecContext(ctx, query, location.UserId, location.Geog)
	if err != nil {
		return 0, err
	}
	ret, err := rows.RowsAffected()
	if err != nil {
		return 0, err
	}
	return ret, nil

}
