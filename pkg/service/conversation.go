package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

var (
	ErrConvoWithSelf = fmt.Errorf("conversation: cannot create conversation with yourself")
)

func NewConversation(convRepo repository.Conversation) *conversation {
	return &conversation{
		convRepo: convRepo,
	}
}

type conversation struct {
	convRepo repository.Conversation
}

func (c *conversation) CreateConversation(conv *domain.Conversation) (string, error) {
	if conv.FromUser.ID == conv.ToUser.ID {
		return "", ErrConvoWithSelf
	}
	//TODO: check if has been match or not
	id, err := c.convRepo.InsertConversation(conv.FromUser.ID, conv.ToUser.ID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrResourceNotFound
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return "", domain.ErrRefNotFound23503
			}
			if pqErr.Code == "23505" {
				return "", domain.ErrUniqueConstraint23505
			}
			return "", pqErr
		}
		return "", err
	}
	return id, nil
}

func (c *conversation) FindConversationById(convoId string) (*domain.Conversation, error) {
	conv, err := c.convRepo.SelectConversationById(convoId)
	//TODO: check if has been reveal, if not, show generic profile_pic and alias
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccesingDB
		}
		return nil, err
	}
	return conv, nil
}

func (c *conversation) GetConversationByUserId(userId string) ([]domain.Conversation, error) {
	convs, err := c.convRepo.SelectConversationByUserId(userId, nil)
	//TODO: check if has been reveal, if not, show generic profile_pic and alias
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccesingDB
		}
		return nil, err
	}
	return convs, nil
}

func (c *conversation) UpdateConvRow(convoId string) error {
	err := c.convRepo.UpdateChatRow(convoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	return nil
}

func (c *conversation) UpdateConvDay(convoId string) error {
	err := c.convRepo.UpdateDayPass(convoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	return nil
}
