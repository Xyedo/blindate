package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/common"
)

func NewAuth(conn *sqlx.DB) *AuthConn {
	return &AuthConn{
		conn: conn,
	}
}

type AuthConn struct {
	conn *sqlx.DB
}

func (a *AuthConn) AddRefreshToken(token string) error {
	query := `INSERT INTO authentications VALUES($1)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := a.conn.ExecContext(ctx, query, token)
	if err != nil {
		return a.wrapError(err)
	}
	ret, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ret == int64(0) {
		return common.ErrResourceNotFound
	}
	return nil

}

func (a *AuthConn) VerifyRefreshToken(token string) error {
	query := `SELECT token FROM authentications WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbToken string
	err := a.conn.QueryRowxContext(ctx, query, token).Scan(&dbToken)
	if err != nil {
		return a.wrapError(err)
	}
	return nil
}

func (a *AuthConn) DeleteRefreshToken(token string) error {
	query := `DELETE FROM authentications WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row, err := a.conn.ExecContext(ctx, query, token)
	if err != nil {
		return a.wrapError(err)
	}
	retInt, err := row.RowsAffected()
	if err != nil {
		return a.wrapError(err)
	}
	if retInt == int64(0) {
		return common.ErrResourceNotFound
	}
	return nil
}
func (AuthConn) wrapError(err error) error {
	var pqErr *pq.Error
	switch {
	case errors.Is(err, context.Canceled):
		return common.WrapError(err, common.ErrTooLongAccessingDB)
	case errors.Is(err, sql.ErrNoRows):
		return common.WrapError(err, common.ErrNotMatchCredential)
	case errors.As(err, &pqErr) && pqErr.Code == "23505":
		return common.WrapErrorWithMsg(err, common.ErrUniqueConstraint23505, "token is already taken, please try again")
	default:
		return err
	}
}
