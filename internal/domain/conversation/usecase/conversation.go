package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/xyedo/blindate/internal/domain/conversation"
	"github.com/xyedo/blindate/internal/domain/conversation/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
	"github.com/xyedo/blindate/pkg/pagination"
)

func New(matchRepo conversation.Repository, userUsecase conversation.UserUsecase, fileUsecase conversation.AttachmentUsecase) *Conversation {
	return &Conversation{
		repo:        matchRepo,
		userUsecase: userUsecase,
	}
}

type Conversation struct {
	repo              conversation.Repository
	userUsecase       conversation.UserUsecase
	attachmentUsecase conversation.AttachmentUsecase
}

var _ conversation.Usecase = &Conversation{}

func (uc *Conversation) IndexConversation(ctx context.Context, requestId string, page, limit int) (entities.ConversationIndex, bool, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return nil, false, err
	}
	defer conn.Release()

	_, err = uc.userUsecase.GetUserDetailByUserId(ctx, conn, requestId)
	if err != nil {
		return nil, false, err
	}

	convos, hasNext, err := uc.repo.FindConversationsByUserId(ctx, conn,
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
		files, err := uc.attachmentUsecase.FindFilesByIds(ctx, conn, fileIds)
		if err != nil {
			return nil, false, err
		}

		var wg sync.WaitGroup
		errs := make([]error, len(files))

		wg.Add(len(files))
		for i := 0; i < len(files); i++ {
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				presignedURL, err := uc.attachmentUsecase.GetAttachmentURL(ctx, files[i].BlobLink, 1*time.Hour)
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

func (uc *Conversation) IndexChatByConversationId(ctx context.Context, payload entities.IndexChatPayload) (entities.Conversation, bool, bool, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return entities.Conversation{}, false, false, err
	}
	defer conn.Release()

	_, err = uc.userUsecase.GetUserDetailByUserId(ctx, conn, payload.RequestId)
	if err != nil {
		return entities.Conversation{}, false, false, err
	}

	conv, hasNext, hasPrev, err := uc.repo.FindChatsByConversationId(ctx, conn, payload)
	if err != nil {
		return entities.Conversation{}, false, false, err
	}

	fileId, ok := conv.Recepient.FileId.Get()
	if !ok {
		return conv, hasNext, hasPrev, nil
	}

	file, err := uc.attachmentUsecase.GetFileById(ctx, conn, fileId)
	if err != nil {
		return entities.Conversation{}, false, false, err
	}

	presignedURL, err := uc.attachmentUsecase.GetAttachmentURL(ctx, file.BlobLink, 1*time.Hour)
	if err != nil {
		return entities.Conversation{}, false, false, err
	}

	conv.Recepient.Url = presignedURL

	return conv, hasNext, hasPrev, nil
}

func (uc *Conversation) CreateConversation(ctx context.Context, conn pg.Querier, payload entities.Conversation) error {
	return uc.repo.CreateConversation(ctx, conn, payload)
}
