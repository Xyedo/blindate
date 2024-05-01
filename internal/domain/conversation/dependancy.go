package conversation

import (
	"context"
	"time"

	attachmentEntities "github.com/xyedo/blindate/internal/domain/attachment/entities"
	userEntities "github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

type UserUsecase interface {
	GetUserDetailByUserId(ctx context.Context, conn pg.Querier, userId string) (userEntities.UserDetail, error)
}

type AttachmentUsecase interface {
	FindFilesByIds(ctx context.Context, conn pg.Querier, ids []string) ([]attachmentEntities.File, error)
	GetFileById(ctx context.Context, conn pg.Querier, id string) (attachmentEntities.File, error)

	GetAttachmentURL(ctx context.Context, key string, expires time.Duration) (string, error)
}
