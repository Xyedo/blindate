package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/common"
	onlineEntities "github.com/xyedo/blindate/pkg/domain/online/entities"
)

func NewOnline(db *sqlx.DB) *OnlineCon {
	return &OnlineCon{
		conn: db,
	}
}

type OnlineCon struct {
	conn *sqlx.DB
}

func (o *OnlineCon) InsertNewOnline(on onlineEntities.DTO) error {
	query := `
	INSERT INTO 
	onlines (user_id,last_online,is_online)
	VALUES ($1,$2,$3)
	RETURNING user_id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var retUserId string
	err := o.conn.GetContext(ctx, &retUserId, query, on.UserId, on.LastOnline, on.IsOnline)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "invalid userId")
			}
			if pqErr.Code == "23505" {
				return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "online already created")
			}
			return err
		}
		if errors.Is(err, context.Canceled) {
			return common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return common.WrapError(err, common.ErrResourceNotFound)
		}
		return err
	}

	return nil
}
func (o *OnlineCon) SelectOnline(userId string) (onlineEntities.DTO, error) {
	query := `
	SELECT
		user_id, last_online, is_online
	FROM onlines
	WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var userOnline onlineEntities.DTO
	err := o.conn.GetContext(ctx, &userOnline, query, userId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return onlineEntities.DTO{}, common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return onlineEntities.DTO{}, common.WrapError(err, common.ErrResourceNotFound)
		}
		return onlineEntities.DTO{}, err
	}
	return userOnline, nil

}
func (o *OnlineCon) UpdateOnline(userId string, online bool) error {
	query := `
	UPDATE onlines SET
		is_online=$1, last_online=COALESCE($2, last_online)
	WHERE user_id=$3
	RETURNING user_id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	if online {
		err := o.conn.GetContext(ctx, &id, query, online, pq.NullTime{}, userId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return common.WrapError(err, common.ErrTooLongAccessingDB)
			}
			if errors.Is(err, sql.ErrNoRows) {
				return common.WrapError(err, common.ErrResourceNotFound)
			}
			return err
		}
	} else {
		err := o.conn.GetContext(ctx, &id, query, online, time.Now(), userId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return common.WrapError(err, common.ErrTooLongAccessingDB)
			}
			if errors.Is(err, sql.ErrNoRows) {
				return common.WrapError(err, common.ErrResourceNotFound)
			}
			return err
		}
	}
	return nil

}
