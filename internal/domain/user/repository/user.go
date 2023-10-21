package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
)

func StoreUser(ctx context.Context, conn *pgx.Conn, id string) error {
	const storeUser = `
	INSERT INTO users(id)
	VALUES($1)
`
	ct, err := conn.Exec(ctx, storeUser, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() != 1 {
		panic("not inserted")
	}

	return nil
}

func DeleteUserById(ctx context.Context, conn *pgx.Conn, id string) error {
	const deleteUserById = `
	UPDATE users SET
		is_deleted = true
	where id = $1
	returning id
`
	var returnedId string
	err := conn.QueryRow(ctx, deleteUserById, id).Scan(&returnedId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(apperror.Payload{
				Error: err,
			})
		}

		return err
	}

	return nil

}
