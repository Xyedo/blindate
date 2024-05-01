package conversation

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/conversation/entities"
)

type Usecase interface {
	IndexChatByConversationId(ctx context.Context, payload entities.IndexChatPayload) (entities.Conversation, bool, bool, error)
	IndexConversation(ctx context.Context, requestId string, page int, limit int) (entities.ConversationIndex, bool, error)
}
