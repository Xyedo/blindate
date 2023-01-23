package service

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain/chat"
	chatEntity "github.com/xyedo/blindate/pkg/domain/chat/entities"
	"github.com/xyedo/blindate/pkg/domain/event"
	"github.com/xyedo/blindate/pkg/domain/match"
	matchEntity "github.com/xyedo/blindate/pkg/domain/match/entities"
)

func NewChat(chatRepo chat.Repository, matchRepo match.Repository) *Chat {
	return &Chat{
		chatRepo:  chatRepo,
		matchRepo: matchRepo,
	}
}

type Chat struct {
	chatRepo  chat.Repository
	matchRepo match.Repository
}

func (c *Chat) CreateNewChat(content *chatEntity.DTO) error {
	matchDAO, err := c.matchRepo.GetMatchById(content.ConversationId)
	if err != nil {
		return err
	}
	if !(content.Author == matchDAO.RequestFrom || content.Author == matchDAO.RequestTo) {
		return common.WrapWithNewError(common.ErrAuthorNotValid, http.StatusForbidden, "author not in this conversation")
	}
	if matchDAO.RequestStatus != string(matchEntity.Accepted) {
		return ErrInvalidMatchStatus
	}
	chatDAO := c.convertToDAO(*content)
	cleanChats := c.sanitizeChat(chatDAO)
	for _, cleanChat := range cleanChats {
		err = c.chatRepo.InsertNewChat(&cleanChat)
		if err != nil {
			return err
		}
	}
	cleanChatDTO := make([]chatEntity.DTO, 0, len(cleanChats))
	for _, cleanChat := range cleanChats {
		cleanChatDTO = append(cleanChatDTO, c.convertToDTO(cleanChat))
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
func (c *Chat) GetMessages(convoId string, filter chat.Filter) ([]chatEntity.DTO, error) {
	chats, err := c.chatRepo.SelectChat(convoId, filter)
	if err != nil {
		return nil, err
	}

	chatsDTO := make([]chatEntity.DTO, 0, len(chats))
	for _, chat := range chats {
		chatsDTO = append(chatsDTO, c.convertToDTO(chat))
	}
	return chatsDTO, nil
}

func (c *Chat) DeleteMessagesById(chatId string) error {
	err := c.chatRepo.DeleteChatById(chatId)
	if err != nil {
		return err
	}
	return nil
}
func (*Chat) sanitizeChat(chat chatEntity.DAO) []chatEntity.DAO {
	chat.Messages = strings.TrimSpace(chat.Messages)
	if chat.Attachment != nil && chat.Messages != "" {
		chatWithoutAttachment := chat
		chatWithoutAttachment.Attachment = nil
		chatWithAttachment := chat
		chatWithAttachment.Messages = ""
		return []chatEntity.DAO{chatWithoutAttachment, chatWithAttachment}
	}
	return []chatEntity.DAO{chat}
}
func (*Chat) convertToDAO(content chatEntity.DTO) chatEntity.DAO {
	chatDAO := chatEntity.DAO{
		Id:             content.Id,
		ConversationId: content.ConversationId,
		Author:         content.Author,
		Messages:       content.Messages,
		SentAt:         content.SentAt,
		Attachment:     content.Attachment,
	}
	if content.ReplyTo != nil {
		chatDAO.ReplyTo = sql.NullString{
			Valid:  true,
			String: *content.ReplyTo,
		}
	}
	if content.SeenAt != nil {
		chatDAO.SeenAt = sql.NullTime{
			Valid: true,
			Time:  *content.SeenAt,
		}
	}

	return chatDAO
}

func (*Chat) convertToDTO(content chatEntity.DAO) chatEntity.DTO {
	chatDomain := chatEntity.DTO{
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
