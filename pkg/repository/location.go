package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/entity"
)

type Location interface {
	InsertNewLocation(location *entity.Location) (int64, error)
	UpdateLocation(location *entity.Location) (int64, error)
	GetLocationByUserId(id string) (*entity.Location, error)
}

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

func (l *location) UpdateLocation(location *entity.Location) (int64, error) {
	query := `
		UPDATE locations SET geog = ST_GeomFromText($1)
		WHERE user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := l.ExecContext(ctx, query, location.Geog, location.UserId)
	if err != nil {
		return 0, nil
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return rows, nil
}

func (l *location) GetLocationByUserId(id string) (*entity.Location, error) {
	query := `
		SELECT ST_AsText(geog) as geog FROM locations WHERE user_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var location entity.Location
	err := l.GetContext(ctx, &location.Geog, query, id)
	if err != nil {
		return nil, err
	}
	location.UserId = id
	return &location, nil

}
