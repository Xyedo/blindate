package attachment

import (
	"context"
	"time"

	"github.com/xyedo/blindate/internal/domain/attachment/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

type Repository interface {
	FindFileByIds(ctx context.Context, conn pg.Querier, ids []string) ([]entities.File, error)
	InsertFile(ctx context.Context, conn pg.Querier, file entities.File) (string, error)
}

type Storage interface {
	DeleteBlob(ctx context.Context, key string) error
	GetPresignedUrl(ctx context.Context, key string, expires time.Duration) (string, error)
	UploadAttachment(ctx context.Context, attachment entities.Attachment) (string, error)
}
