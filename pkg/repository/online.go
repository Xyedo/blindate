package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
)

// TODO: TEST Online Repository
type Online interface {
	InsertNewOnline(on *domain.Online) error
	UpdateOnline(userId string, online bool) error
	SelectOnline(userId string) (*domain.Online, error)
}

func NewOnline(db *sqlx.DB) *online {
	return &online{
		db,
	}
}

type online struct {
	*sqlx.DB
}

func (o *online) InsertNewOnline(on *domain.Online) error {
	query := `
	INSERT INTO 
	onlines (user_id,last_online,is_online)
	VALUES ($1,$2,$3)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := o.ExecContext(ctx, query, on.UserId, on.LastOnline, on.IsOnline)
	if err != nil {
		return err
	}

	return nil
}
func (o *online) SelectOnline(userId string) (*domain.Online, error) {
	query := `
	SELECT
		user_id, last_online, is_online
	FROM onlines
	WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var userOnline domain.Online
	err := o.GetContext(ctx, &userOnline, query, userId)
	if err != nil {
		return nil, err
	}
	return &userOnline, nil

}
func (o *online) UpdateOnline(userId string, online bool) error {
	query := `
	UPDATE onlines SET
		is_online=$1, last_online=COALESCE($2, last_online)
	WHERE user_id=$3
	RETURNING user_id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	if online {
		err := o.GetContext(ctx, &id, query, online, pq.NullTime{}, userId)
		if err != nil {
			return err
		}
	} else {
		err := o.GetContext(ctx, &id, query, online, time.Now(), userId)
		if err != nil {
			return err
		}
	}
	return nil

}
