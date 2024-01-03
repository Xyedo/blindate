package usecase

import (
	"context"
	"sync"
	"time"

	attachmentRepo "github.com/xyedo/blindate/internal/domain/attachment/repository"
	"github.com/xyedo/blindate/internal/domain/attachment/s3"
	"github.com/xyedo/blindate/internal/domain/conversation/entities"
	"github.com/xyedo/blindate/internal/domain/conversation/repository"
	userRepo "github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
	"github.com/xyedo/blindate/pkg/pagination"
)

func IndexConversation(ctx context.Context, requestId string, page, limit int) (entities.ConversationIndex, bool, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return nil, false, err
	}
	defer conn.Release()

	_, err = userRepo.GetUserDetailById(ctx, conn, requestId)
	if err != nil {
		return nil, false, err
	}

	convos, hasNext, err := repository.FindConversationsByUserId(ctx, conn,
		requestId,
		pagination.Pagination{
			Page:  page,
			Limit: limit,
		},
	)
	if err != nil {
		return nil, false, err
	}

	fileIds, fileIdToConvosIdx := convos.ToFileIds()
	if len(fileIds) > 0 {
		files, err := attachmentRepo.GetFileByIds(ctx, conn, fileIds)
		if err != nil {
			return nil, false, err
		}

		var wg sync.WaitGroup
		errs := make([]error, len(files))

		wg.Add(len(files))
		for i := 0; i < len(files); i++ {
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				presignedURL, err := s3.Manager.GetPresignedUrl(ctx, files[i].BlobLink, 1*time.Hour)
				if err != nil {
					errs[i] = err
					return
				}

				if idx, ok := fileIdToConvosIdx[files[i].Id]; ok {
					convos[idx].Recepient.Url = presignedURL
				}
			}(i, &wg)
		}
		wg.Wait()

		for _, err := range errs {
			if err != nil {
				return nil, false, err
			}
		}
	}

	return convos, hasNext, nil
}

func IndexChatByConversationId(ctx context.Context, payload entities.IndexChatPayload) (entities.Conversation, bool, bool, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return entities.Conversation{}, false, false, err
	}
	defer conn.Release()

	_, err = userRepo.GetUserDetailById(ctx, conn, payload.RequestId)
	if err != nil {
		return entities.Conversation{}, false, false, err
	}

	conv, hasNext, hasPrev, err := repository.FindChatsByConversationId(ctx, conn, payload)
	if err != nil {
		return entities.Conversation{}, false, false, err
	}

	fileId, ok := conv.Recepient.FileId.Get()
	if !ok {
		return conv, hasNext, hasPrev, nil
	}

	files, err := attachmentRepo.GetFileByIds(ctx, conn, []string{fileId})
	if err != nil {
		return entities.Conversation{}, false, false, err
	}

	presignedURL, err := s3.Manager.GetPresignedUrl(ctx, files[0].BlobLink, 1*time.Hour)
	if err != nil {
		return entities.Conversation{}, false, false, err
	}

	conv.Recepient.Url = presignedURL

	return conv, hasNext, hasPrev, nil
}
