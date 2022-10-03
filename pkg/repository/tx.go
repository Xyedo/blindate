package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func execTx(conn *sqlx.DB, ctx context.Context, cb func(q *sqlx.DB) error) error {
	tx, err := conn.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
	if err != nil {
		return err
	}
	err = cb(conn)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
