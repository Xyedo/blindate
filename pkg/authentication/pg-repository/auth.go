package pgrepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/authentication"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
)

func NewAuth(conn *sqlx.DB) authentication.Repository {
	return &AuthConn{
		conn: conn,
	}
}

type AuthConn struct {
	conn *sqlx.DB
}

func (a *AuthConn) AddRefreshToken(token string) error {
	query := `INSERT INTO authentications VALUES($1) RETURNING token`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnedToken string
	err := a.conn.GetContext(ctx, &returnedToken, query, token)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.Is(err, context.Canceled):
			return apperror.Timeout(apperror.Payload{Error: err})
		case errors.Is(err, sql.ErrNoRows):
			return apperror.Unauthorized(apperror.Payload{Error: err})
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			return apperror.Conflicted(apperror.Payload{Error: err, Message: "token is already taken, please try again"})
		default:
			return err
		}
	}

	return nil

}

func (a *AuthConn) VerifyRefreshToken(token string) error {
	query := `SELECT token FROM authentications WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbToken string
	err := a.conn.GetContext(ctx, &dbToken, query, token)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return apperror.Timeout(apperror.Payload{Error: err})
		case errors.Is(err, sql.ErrNoRows):
			return apperror.Unauthorized(apperror.Payload{Error: err})
		default:
			return err
		}
	}
	return nil
}

func (a *AuthConn) DeleteRefreshToken(token string) error {
	query := `DELETE FROM authentications WHERE token = $1 RETURNING token`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returningToken string
	err := a.conn.GetContext(ctx, &returningToken, query, token)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			return apperror.Timeout(apperror.Payload{Error: err})
		case errors.Is(err, sql.ErrNoRows):
			return apperror.Unauthorized(apperror.Payload{Error: err})
		default:
			return err
		}
	}
	return nil
}
