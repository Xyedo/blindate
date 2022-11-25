package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
	"github.com/xyedo/blindate/pkg/repository"
)

var (
	ErrRefMediaType   = fmt.Errorf("%w:invalid media types", domain.ErrRefNotFound23503)
	ErrRefConvoID     = fmt.Errorf("%w:invalid convoId", domain.ErrRefNotFound23503)
	ErrRefReplyTo     = fmt.Errorf("%w:invalid reply_to", domain.ErrRefNotFound23503)
	ErrAuthorNotValid = errors.New("author not in the conversation")
)

func NewChat(chatRepo repository.Chat, matchRepo repository.Match) *chat {
	return &chat{
		chatRepo:  chatRepo,
		matchRepo: matchRepo,
	}
}

type chat struct {
	chatRepo  repository.Chat
	matchRepo repository.Match
}

func (c *chat) CreateNewChat(content *domain.Chat) error {

	chatEntity := c.convertToEntity(content)
	matchEntity, err := c.matchRepo.GetMatchById(chatEntity.ConversationId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRefConvoID
		}
		return err
	}
	if !(chatEntity.Author == matchEntity.RequestFrom || chatEntity.Author == matchEntity.RequestTo) {
		return ErrAuthorNotValid
	}
	if matchEntity.RequestStatus != string(domain.Accepted) {
		return ErrNotYetAccepted
	}
	cleanChats := c.sanitizeChat(chatEntity)
	for _, cleanChat := range cleanChats {
		err = c.chatRepo.InsertNewChat(cleanChat)
		if err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				if pqErr.Code == "23503" {
					if strings.Contains(pqErr.Constraint, "media_type") {
						return ErrRefMediaType
					}
					if strings.Contains(pqErr.Constraint, "conversation_id") {
						return ErrRefConvoID
					}
					if strings.Contains(pqErr.Constraint, "author") {
						return ErrAuthorNotValid
					}
					if strings.Contains(pqErr.Constraint, "reply_to") {
						return ErrRefReplyTo
					}
				}
				return pqErr
			}
			return err
		}
	}
	return nil

}

func (c *chat) GetMessages(convoId string, filter entity.ChatFilter) ([]domain.Chat, error) {
	chats, err := c.chatRepo.SelectChat(convoId, filter)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccesingDB
		}
		return nil, err
	}

	chatDomain := make([]domain.Chat, 0, len(chats))
	for _, chat := range chats {
		chatDomain = append(chatDomain, *c.convertToDomain(&chat))
	}
	return chatDomain, nil
}

func (c *chat) DeleteMessagesById(chatId string) error {
	err := c.chatRepo.DeleteChatById(chatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrRefNotFound23503
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	return nil
}
func (*chat) sanitizeChat(chat *entity.Chat) []*entity.Chat {
	chat.Messages = strings.TrimSpace(chat.Messages)
	if chat.Attachment != nil && chat.Messages != "" {
		chatWAttach := *chat
		chatWAttach.Messages = ""
		chatWoAttach := *chat
		chatWoAttach.Attachment = nil
		return []*entity.Chat{&chatWAttach, &chatWoAttach}
	}
	return []*entity.Chat{chat}
}
func (*chat) convertToEntity(content *domain.Chat) *entity.Chat {
	chatEntity := &entity.Chat{
		Id:             content.Id,
		ConversationId: content.ConversationId,
		Author:         content.Author,
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
		Author:         content.Author,
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
