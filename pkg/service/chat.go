package service

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/domain/entity"
	"github.com/xyedo/blindate/pkg/event"
	"github.com/xyedo/blindate/pkg/repository"
)

func NewChat(chatRepo repository.Chat, matchRepo repository.Match) *Chat {
	return &Chat{
		chatRepo:  chatRepo,
		matchRepo: matchRepo,
	}
}

type Chat struct {
	chatRepo  repository.Chat
	matchRepo repository.Match
}

func (c *Chat) CreateNewChat(content *domain.Chat) error {
	chatEntity := c.convertToEntity(content)
	matchEntity, err := c.matchRepo.GetMatchById(chatEntity.ConversationId)
	if err != nil {
		return err
	}
	if !(chatEntity.Author == matchEntity.RequestFrom || chatEntity.Author == matchEntity.RequestTo) {
		return common.WrapWithNewError(common.ErrAuthorNotValid, http.StatusForbidden, "author not in this conversation")
	}
	if matchEntity.RequestStatus != string(domain.Accepted) {
		return ErrInvalidMatchStatus
	}
	cleanChats := c.sanitizeChat(chatEntity)
	for _, cleanChat := range cleanChats {
		err = c.chatRepo.InsertNewChat(cleanChat)
		if err != nil {
			return err
		}
	}
	cleanChatDTO := make([]domain.Chat, 0, len(cleanChats))
	for i := range cleanChats {
		cleanChatDTO = append(cleanChatDTO, *c.convertToDomain(cleanChats[i]))
	}
	event.ChatCreated.Trigger(event.ChatCreatedPayload{
		Chat:   cleanChatDTO,
		ConvId: content.ConversationId,
	})
	return nil
}
func (c *Chat) UpdateSeenChat(convId, userId string) error {
	matchEntity, err := c.matchRepo.GetMatchById(convId)
	if err != nil {
		return err
	}
	if !(matchEntity.RequestFrom == userId || matchEntity.RequestTo == userId) {
		return common.WrapWithNewError(common.ErrAuthorNotValid, http.StatusForbidden, "users not in this conversation")
	}
	changedChatIds, err := c.chatRepo.UpdateSeenChat(convId, userId)
	if err != nil {
		return err
	}

	event.ChatSeen.Trigger(event.ChatSeenPayload{
		RequestFrom: matchEntity.RequestFrom,
		RequestTo:   matchEntity.RequestTo,
		SeenChatIds: changedChatIds,
	})
	return nil
}
func (c *Chat) GetMessages(convoId string, filter entity.ChatFilter) ([]domain.Chat, error) {
	chats, err := c.chatRepo.SelectChat(convoId, filter)
	if err != nil {
		return nil, err
	}

	chatDomain := make([]domain.Chat, 0, len(chats))
	for _, chat := range chats {
		chatDomain = append(chatDomain, *c.convertToDomain(&chat))
	}
	return chatDomain, nil
}

func (c *Chat) DeleteMessagesById(chatId string) error {
	err := c.chatRepo.DeleteChatById(chatId)
	if err != nil {
		return err
	}
	return nil
}
func (*Chat) sanitizeChat(chat *entity.Chat) []*entity.Chat {
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
func (*Chat) convertToEntity(content *domain.Chat) *entity.Chat {
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

func (*Chat) convertToDomain(content *entity.Chat) *domain.Chat {
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
