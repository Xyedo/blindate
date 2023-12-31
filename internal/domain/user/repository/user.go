package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func StoreUser(ctx context.Context, conn pg.Querier, id string) error {
	const storeUser = `
	INSERT INTO account(id)
	VALUES($1)
	ON CONFLICT(id)
	DO UPDATE SET
	is_deleted = false
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

func GetUserById(ctx context.Context, conn pg.Querier, id string, opts ...entities.GetUserOption) (entities.User, error) {
	const storeUser = `
	SELECT 
		id, 
		is_deleted 
	FROM account 
	WHERE id = $1
`
	query := storeUser

	if !(len(opts) > 0 && opts[0].WithDeleted) {
		query += " AND is_deleted = false"
	}

	if len(opts) > 0 && opts[0].PessimisticLocking {
		query += "\n FOR UPDATE"
	}

	var returnedUser entities.User
	err := conn.QueryRow(ctx, query, id).Scan(
		&returnedUser.Id,
		&returnedUser.IsDeleted,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.User{}, apperror.NotFound(apperror.Payload{
				Status: entities.ErrCodeUserNotFound,
				Error:  err,
			})
		}

		return entities.User{}, err
	}
	return returnedUser, nil
}

func SoftDeleteUserById(ctx context.Context, conn pg.Querier, id string) error {
	const deleteUserById = `
	UPDATE account SET
		is_deleted = true
	where id = $1
	returning id
`
	var returnedId string
	err := conn.QueryRow(ctx, deleteUserById, id).Scan(&returnedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperror.NotFound(apperror.Payload{
				Status: entities.ErrCodeUserNotFound,
				Error:  err,
			})
		}

		return err
	}

	return nil

}
