package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Auth interface {
	AddRefreshToken(token string) (int64, error)
	VerifyRefreshToken(token string) error
	DeleteRefreshToken(token string) (int64, error)
}

func NewAuth(conn *sqlx.DB) *authConn {
	return &authConn{
		conn: conn,
	}
}

type authConn struct {
	conn *sqlx.DB
}

func (a *authConn) AddRefreshToken(token string) (int64, error) {
	query := `INSERT INTO authentications VALUES($1)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := a.conn.ExecContext(ctx, query, token)
	if err != nil {
		return 0, err
	}
	ret, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return ret, nil

}

func (a *authConn) VerifyRefreshToken(token string) error {
	query := `SELECT token FROM authentications WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbToken string
	err := a.conn.QueryRowxContext(ctx, query, token).Scan(&dbToken)
	if err != nil {
		return err
	}
	return nil
}

func (a *authConn) DeleteRefreshToken(token string) (int64, error) {
	query := `DELETE FROM authentications WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row, err := a.conn.ExecContext(ctx, query, token)
	if err != nil {
		return 0, err
	}
	retInt, err := row.RowsAffected()
	if err != nil {
		return 0, err
	}
	return retInt, nil
}
