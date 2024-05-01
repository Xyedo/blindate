package user

import (
	"context"
	"time"

	"github.com/xyedo/blindate/internal/domain/attachment/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

type AttachmentUsecase interface {
	FindFilesByIds(ctx context.Context, conn pg.Querier, ids []string) ([]entities.File, error)
	InsertFile(ctx context.Context, conn pg.Querier, file entities.File) (string, error)

	UploadAttachment(ctx context.Context, attachment entities.Attachment) (string, error)
	GetAttachmentURL(ctx context.Context, key string, expires time.Duration) (string, error)
}
