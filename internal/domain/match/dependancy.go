package match

import (
	"context"

	conversationEntities "github.com/xyedo/blindate/internal/domain/conversation/entities"
	userEntities "github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

type UserUsecase interface {
	GetUserDetailByUserId(ctx context.Context, conn pg.Querier, userId string) (userEntities.UserDetail, error)
	GetUserDetails(ctx context.Context, conn pg.Querier, userIds []string) (userEntities.UserDetails, error)

	FindNonMatchClosestUserIds(ctx context.Context, conn pg.Querier, payload userEntities.FindClosestUser) ([]string, error)
}

type ConversationUsecase interface {
	CreateConversation(ctx context.Context, conn pg.Querier, payload conversationEntities.Conversation) error
}
