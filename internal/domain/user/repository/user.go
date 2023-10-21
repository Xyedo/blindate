package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const storeUser = `
	INSERT INTO users(id)
	VALUES($1)
`

func StoreUser(ctx context.Context, conn *pgx.Conn, id string) error {
	ct, err := conn.Exec(ctx, storeUser, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() != 1 {
		panic("not inserted")
	}

	return nil
}
