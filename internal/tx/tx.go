package tx

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// ExecGeneric need to be to use in your class method
func ExecGeneric(conn *sqlx.DB, ctx context.Context, cb func(tx *sqlx.Tx) error, option *sql.TxOptions) error {
	tx, err := conn.BeginTxx(ctx, option)
	if err != nil {
		return err
	}
	err = cb(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, rb err: %w", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
