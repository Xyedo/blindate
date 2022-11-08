package service

import (
	"database/sql"

	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
)

type chatRepo interface {
	InsertNewChat(content *entity.Chat) error
	GetChats(convoId string, limit int) ([]entity.Chat, error)
	DeleteChat(chatId string) (int64, error)
}

type attachmentSvc interface {
	UploadBlob(file []byte, contentType string) (string, error)
}

func NewChat(chatRepo chatRepo, chatAttachment attachmentSvc) *chat {
	return &chat{
		chatRepo:      chatRepo,
		attachmentSvc: chatAttachment,
	}
}

type chat struct {
	chatRepo      chatRepo
	attachmentSvc attachmentSvc
}

func (c *chat) CreateNewChat(content *domain.Chat) error {
	chatEntity := c.convertToEntity(content)
	err := c.chatRepo.InsertNewChat(chatEntity)
	if err != nil {
		return err
	}
	return nil
}

func (c *chat) CreateNewChatwithMedia(content *domain.Chat, file []byte, contentType string) error {
	//TODO: if contentType == audio and have messages, send it two chat, one only messsages, one only audio
	chatEntity := c.convertToEntity(content)
	ref, err := c.attachmentSvc.UploadBlob(file, contentType)
	if err != nil {
		return err
	}
	chatEntity.Attachment = &entity.Attachment{
		ChatId:    chatEntity.Id,
		BlobLink:  ref,
		MediaType: contentType,
	}
	err = c.chatRepo.InsertNewChat(chatEntity)
	if err != nil {
		return err
	}
	return nil
}

func (*chat) convertToEntity(content *domain.Chat) *entity.Chat {
	chatEntity := &entity.Chat{
		Id:             content.Id,
		ConversationId: content.ConversationId,
		Messages:       content.Messages,
		SentAt:         content.SentAt,
	}
	if content.ReplyTo != nil {
		chatEntity.ReplyTo = sql.NullString{
			Valid:  true,
			String: *content.ReplyTo,
		}
	}
	return chatEntity
}
