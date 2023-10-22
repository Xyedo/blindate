package usecase

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func RegisterUser(ctx context.Context, id string) error {
	pool, err := pg.GetPool(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	return repository.StoreUser(ctx, pool, id)
}

func DeleteUser(ctx context.Context, id string) error {
	pool, err := pg.GetPool(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	return repository.DeleteUserById(ctx, pool, id)
}
