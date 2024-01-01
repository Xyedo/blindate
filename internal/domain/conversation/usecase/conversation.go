package usecase

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/conversation/entities"
)

func IndexConversation(ctx context.Context, requestId string, page, limit int) (entities.Conversations, error) {
	panic("unimplemented")
}
