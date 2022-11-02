package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func execGeneric(conn *sqlx.DB, ctx context.Context, cb func(q *sqlx.DB) error, option *sql.TxOptions) error {
	tx, err := conn.BeginTxx(ctx, option)
	if err != nil {
		return err
	}
	err = cb(conn)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, rb err: %w", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
