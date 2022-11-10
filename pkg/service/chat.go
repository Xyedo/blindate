package service

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
)

type chatRepo interface {
	InsertNewChat(content *entity.Chat) error
	SelectChat(convoId string, filter entity.ChatFilter) ([]entity.Chat, error)
	DeleteChat(chatId string) (int64, error)
}

func NewChat(chatRepo chatRepo) *chat {
	return &chat{
		chatRepo: chatRepo,
	}
}

type chat struct {
	chatRepo chatRepo
}

func (c *chat) CreateNewChat(content *domain.Chat) error {
	//TODO: check if has been match or not

	chatEntity := c.convertToEntity(content)
	err := c.chatRepo.InsertNewChat(chatEntity)
	if err != nil {
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

func (c *chat) GetMessages(convoId string, filter entity.ChatFilter) ([]domain.Chat, error) {
	chats, err := c.chatRepo.SelectChat(convoId, filter)
	if err != nil {
		//TODO: make error handling more good
		return nil, err
	}

	chatDomain := make([]domain.Chat, 0, len(chats))
	for _, chat := range chats {
		chatDomain = append(chatDomain, *c.convertToDomain(&chat))
	}
	return chatDomain, nil
}

func (c *chat) DeleteMessages(chatId string) {
	//TODO: make error handling better
	c.chatRepo.DeleteChat(chatId)
}
func (*chat) convertToEntity(content *domain.Chat) *entity.Chat {
	chatEntity := &entity.Chat{
		Id:             content.Id,
		ConversationId: content.ConversationId,
		Messages:       content.Messages,
		SentAt:         content.SentAt,
		Attachment:     content.Attachment,
	}
	if content.ReplyTo != nil {
		chatEntity.ReplyTo = sql.NullString{
			Valid:  true,
			String: *content.ReplyTo,
		}
	}
	if content.SeenAt != nil {
		chatEntity.SeenAt = sql.NullTime{
			Valid: true,
			Time:  *content.SeenAt,
		}
	}

	return chatEntity
}

func (*chat) convertToDomain(content *entity.Chat) *domain.Chat {
	chatDomain := &domain.Chat{
		Id:             content.Id,
		ConversationId: content.ConversationId,
		Messages:       content.Messages,
		SentAt:         content.SentAt,
		Attachment:     content.Attachment,
	}
	if content.ReplyTo.Valid {
		chatDomain.ReplyTo = &content.ReplyTo.String
	}
	if content.SeenAt.Valid {
		chatDomain.SeenAt = &content.SeenAt.Time
	}
	return chatDomain
}
