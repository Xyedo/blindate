package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/domain/user"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func New(userRepo user.Repository, fileUsecase user.AttachmentUsecase) *User {
	return &User{
		repo:              userRepo,
		attachmentUsecase: fileUsecase,
	}
}

var _ user.Usecase = &User{}

type User struct {
	repo              user.Repository
	attachmentUsecase user.AttachmentUsecase
}

func (uc *User) RegisterUser(ctx context.Context, id string) error {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	return uc.repo.StoreUser(ctx, conn, id)
}

func (uc *User) DeleteUser(ctx context.Context, id string) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		_, err := uc.repo.GetUserById(ctx, tx, id, entities.GetUserOption{
			PessimisticLocking: true,
		})
		if err != nil {
			return err
		}
		return uc.repo.SoftDeleteUserById(ctx, tx, id)
	})

}

func (uc *User) FindNonMatchClosestUserIds(ctx context.Context, conn pg.Querier, payload entities.FindClosestUser) ([]string, error) {
	return uc.repo.FindNonMatchClosestUser(ctx, conn, payload)
}
