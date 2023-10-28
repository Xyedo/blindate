package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func RegisterUser(ctx context.Context, id string) error {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return repository.StoreUser(ctx, conn, id)
}

func DeleteUser(ctx context.Context, id string) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		_, err := repository.GetUserById(ctx, tx, id, entities.GetUserOption{
			PessimisticLocking: true,
		})
		if err != nil {
			return err
		}
		return repository.SoftDeleteUserById(ctx, tx, id)

	})

}
