package usecase

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func RegisterUser(ctx context.Context, id string) error {
	return repository.StoreUser(ctx, pg.Get(), id)
}

func DeleteUser(ctx context.Context, id string) error {
	return repository.DeleteUserById(ctx, pg.Get(), id)
}
