package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
)

var (
	ErrConvoWithSelf = fmt.Errorf("conversation: cannot create conversation with yourself")
)

type conversationRepo interface {
	InsertConversation(convo *domain.Conversation) error
}

func NewConversation(convRepo conversationRepo) *conversation {
	return &conversation{
		convRepo: convRepo,
	}
}

type conversation struct {
	convRepo conversationRepo
}

func (c *conversation) CreateConversation(conv *domain.Conversation) error {
	if conv.FromId == conv.ToId {
		return ErrConvoWithSelf
	}
	//TODO: check if has been match or not
	err := c.convRepo.InsertConversation(conv)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrResourceNotFound
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return domain.ErrRefNotFound23503
			}
			if pqErr.Code == "23505" {
				return domain.ErrUniqueConstraint23505
			}
			return pqErr
		}
		return err
	}
	return nil
}
