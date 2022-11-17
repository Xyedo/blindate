package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
)

type LocationRepo interface {
	InsertNewLocation(location *entity.Location) (int64, error)
	UpdateLocation(location *entity.Location) (int64, error)
	GetLocationByUserId(id string) (*entity.Location, error)
	GetClosestUser(userId string, limit int) ([]domain.User, error)
}

func NewLocation(db *sqlx.DB) *location {
	return &location{
		conn: db,
	}
}

type location struct {
	conn *sqlx.DB
}

func (l *location) InsertNewLocation(location *entity.Location) (int64, error) {
	query := `
		INSERT INTO locations(user_id, geog, created_at, updated_at)
		VALUES($1, ST_GeomFromText($2), $3, $3)`
	now := time.Now()
	args := []any{location.UserId, location.Geog, now}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := l.conn.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	ret, err := rows.RowsAffected()
	if err != nil {
		return 0, err
	}
	location.CreatedAt = now
	location.UpdatedAt = now
	return ret, nil
}

func (l *location) UpdateLocation(location *entity.Location) (int64, error) {
	query := `
		UPDATE locations SET geog = ST_GeomFromText($1), updated_at = $2
		WHERE user_id = $3`
	args := []any{location.Geog, time.Now(), location.UserId}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := l.conn.ExecContext(ctx, query, args...)
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
		SELECT 
			user_id,
			ST_AsText(geog) as geog,
			created_at, 
			updated_at 
		FROM locations 
		WHERE user_id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var location entity.Location
	err := l.conn.GetContext(ctx, &location, query, id)
	if err != nil {
		return nil, err
	}
	return &location, nil

}

func (l *location) GetClosestUser(geom string, limit int) ([]domain.User, error) {
	query := `
		SELECT 
			users.*
		FROM locations
		JOIN users
			ON users.id = locations.user_id
		ORDER BY locations.geog <-> ST_GeomFromText($1)
		LIMIT $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	users := make([]domain.User, 0, limit)
	err := l.conn.SelectContext(ctx, &users, query, geom, limit)
	if err != nil {
		return []domain.User{}, err
	}
	return users, nil
}
