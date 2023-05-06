package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	"github.com/xyedo/blindate/pkg/domain/online"
	onlineEntities "github.com/xyedo/blindate/pkg/domain/online/entities"
)

func New(db *sqlx.DB) online.Repository {
	return &onlineConn{
		conn: db,
	}
}

type onlineConn struct {
	conn *sqlx.DB
}

// Insert implements online.Repository
func (o *onlineConn) Insert(online onlineEntities.Online) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returningUserId string
	err := o.conn.GetContext(
		ctx,
		&returningUserId,
		insertOnline,
		online.UserId,
		online.LastOnline,
		online.IsOnline,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}

		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(apperror.Payload{Error: err})
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if strings.Contains(pqErr.Constraint, "user_id") {
				if pqErr.Code == "23503" {
					return apperror.UnprocessableEntity(apperror.PayloadMap{
						Error: err,
						ErrorMap: map[string]string{
							"user_id": "value not found",
						},
					})
				}
				if pqErr.Code == "23505" {
					return apperror.Conflicted(apperror.Payload{
						Error:   err,
						Message: "online already created",
					})
				}
			}
		}
		return err
	}

	return nil
}

// SelectByUserId implements online.Repository
func (o *onlineConn) SelectByUserId(id string) (onlineEntities.Online, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var online onlineEntities.Online
	err := o.conn.GetContext(ctx, &online, selectOnline, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return onlineEntities.Online{}, apperror.Timeout(apperror.Payload{Error: err})
		}

		if errors.Is(err, sql.ErrNoRows) {
			return onlineEntities.Online{}, apperror.NotFound(apperror.Payload{Error: err})
		}
		return onlineEntities.Online{}, err
	}

	return online, err
}

// Update implements online.Repository
func (o *onlineConn) Update(userId string, online bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var currentTime pq.NullTime
	if !online {
		currentTime.Valid = true
		currentTime.Time = time.Now()
	}

	var returningUserId string
	err := o.conn.GetContext(
		ctx,
		&returningUserId,
		updateOnline,
		online,
		currentTime,
		userId,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}

		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(apperror.Payload{Error: err})
		}
		return err
	}

	return nil
}
