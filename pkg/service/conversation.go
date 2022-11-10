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
	SelectConversationById(convoId string) (*domain.Conversation, error)
	SelectConversationByUserId(UserId string) ([]domain.Conversation, error)
	UpdateDayPass(convoId string) error
	UpdateChatRow(convoId string) error
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
	if conv.FromUser.ID == conv.ToUser.ID {
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
	convs, err := c.convRepo.SelectConversationByUserId(userId)
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
