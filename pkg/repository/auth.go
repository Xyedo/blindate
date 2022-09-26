package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewAuth(conn *sqlx.DB) *authConn {
	return &authConn{
		conn,
	}
}

type authConn struct {
	*sqlx.DB
}

func (a *authConn) AddRefreshToken(token string) (int64, error) {
	query := `INSERT INTO authentications VALUES($1)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := a.DB.ExecContext(ctx, query, token)
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
	err := a.DB.QueryRowxContext(ctx, query, token).Scan(&dbToken)
	if err != nil {
		return err
	}
	return nil
}

func (a *authConn) DeleteRefreshToken(token string) (int64, error) {
	query := `DELETE FROM authentications WHERE token = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row, err := a.DB.ExecContext(ctx, query, token)
	if err != nil {
		return 0, err
	}
	retInt, err := row.RowsAffected()
	if err != nil {
		return 0, err
	}
	return retInt, nil
}
