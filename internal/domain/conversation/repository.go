package conversation

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/conversation/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
	"github.com/xyedo/blindate/pkg/pagination"
)

type Repository interface {
	CreateConversation(ctx context.Context, conn pg.Querier, payload entities.Conversation) error
	FindChatsByConversationId(ctx context.Context, conn pg.Querier, payload entities.IndexChatPayload) (entities.Conversation, bool, bool, error)
	FindConversationsByUserId(ctx context.Context, conn pg.Querier, userId string, pagination pagination.Pagination) (entities.ConversationIndex, bool, error)
}
