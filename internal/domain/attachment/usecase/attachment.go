package usecase

import (
	"context"
	"time"

	"github.com/xyedo/blindate/internal/domain/attachment"
	"github.com/xyedo/blindate/internal/domain/attachment/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func New(attachmentRepo attachment.Repository, storage attachment.Storage) *Attachment {
	return &Attachment{
		repo:    attachmentRepo,
		storage: storage,
	}

}

type Attachment struct {
	repo    attachment.Repository
	storage attachment.Storage
}

func (uc *Attachment) FindFilesByIds(ctx context.Context, conn pg.Querier, ids []string) ([]entities.File, error) {
	return uc.repo.FindFileByIds(ctx, conn, ids)
}

func (uc *Attachment) GetFileById(ctx context.Context, conn pg.Querier, id string) (entities.File, error) {
	files, err := uc.repo.FindFileByIds(ctx, conn, []string{id})
	if err != nil {
		return entities.File{}, err
	}

	if len(files) == 0 {
		return entities.File{}, entities.ErrFileNotFound
	}

	return files[0], nil
}

func (uc *Attachment) InsertFile(ctx context.Context, conn pg.Querier, file entities.File) (string, error) {
	return uc.repo.InsertFile(ctx, conn, file)
}

func (uc *Attachment) UploadAttachment(ctx context.Context, attachment entities.Attachment) (string, error) {
	return uc.storage.UploadAttachment(ctx, attachment)
}

func (uc *Attachment) GetAttachmentURL(ctx context.Context, key string, expires time.Duration) (string, error) {
	return uc.storage.GetPresignedUrl(ctx, key, expires)
}

func (uc *Attachment) DeleteAttachment(ctx context.Context, key string) error {
	return uc.storage.DeleteBlob(ctx, key)
}
