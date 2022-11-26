package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

type chatSvc interface {
	CreateNewChat(content *domain.Chat) error
	GetMessages(convoId string, filter entity.ChatFilter) ([]domain.Chat, error)
	DeleteMessagesById(chatId string) error
}

func NewChat(chatSvc chatSvc, attachSvc attachmentManager) *chat {
	return &chat{
		chatSvc:   chatSvc,
		attachSvc: attachSvc,
	}
}

type chat struct {
	chatSvc   chatSvc
	attachSvc attachmentManager
}

func (chat *chat) postChatHandler(c *gin.Context) {
	var input struct {
		Message string  `json:"message" binding:"required,max=4096"`
		ReplyTo *string `json:"replyTo" binding:"omitempty,uuid"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		if jsonErr := jsonBindingErrResp(err, c, map[string]string{
			"Message": "must not empty and max characters is 4096",
			"ReplyTo": "if specified, must be valid uuid",
		}); jsonErr != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	convoId := c.GetString("convId")
	userId := c.GetString("userId")
	dtoChat := domain.Chat{
		ConversationId: convoId,
		Author:         userId,
		Messages:       input.Message,
		ReplyTo:        input.ReplyTo,
		SentAt:         time.Now(),
	}
	if err := chat.chatSvc.CreateNewChat(&dtoChat); err != nil {
		switch {
		case errors.Is(err, service.ErrRefMediaType):
			errUnprocessableEntityResp(c, "invalid media types")
		case errors.Is(err, service.ErrRefConvoID):
			errUnprocessableEntityResp(c, "conversationId is invalid")
		case errors.Is(err, service.ErrAuthorNotValid):
			errUnprocessableEntityResp(c, "author is invalid user")
		case errors.Is(err, service.ErrRefReplyTo):
			errUnprocessableEntityResp(c, "replyTo is invalid chatId")
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "chat media uploaded",
		"data": gin.H{
			"chat": gin.H{
				"id": dtoChat.Id,
			},
		},
	})

}

func (chat *chat) postChatMediaHandler(c *gin.Context) {
	var validAudioTypes = []string{
		"application/ogg",
		"audio/mpeg",
	}
	key, mediaType := uploadFile(c, chat.attachSvc, validAudioTypes, "chat-attachment")
	if key == "" {
		return
	}
	var input struct {
		ReplyTo *string `form:"replyTo" binding:"omitempty,uuid"`
	}

	userId := c.GetString("userId")
	convoId := c.GetString("convId")
	if err := c.ShouldBindQuery(&input); err != nil {
		if errMap := util.ReadValidationErr(err, map[string]string{
			"ReplyTo": "if provided, it should be in uuid format",
		}); errMap != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "please refer to the documentation",
				"errors":  errMap,
			})
			return
		}
		errServerResp(c, err)
		return
	}
	dtoChat := domain.Chat{
		ConversationId: convoId,
		Author:         userId,
		Messages:       "",
		ReplyTo:        input.ReplyTo,
		SentAt:         time.Now(),
		Attachment: &domain.ChatAttachment{
			BlobLink:  key,
			MediaType: mediaType,
		},
	}
	if err := chat.chatSvc.CreateNewChat(&dtoChat); err != nil {
		switch {
		case errors.Is(err, service.ErrRefMediaType):
			errUnprocessableEntityResp(c, "invalid media types")
		case errors.Is(err, service.ErrRefConvoID):
			errUnprocessableEntityResp(c, "conversationId is invalid")
		case errors.Is(err, service.ErrAuthorNotValid):
			errUnprocessableEntityResp(c, "author is invalid user")
		case errors.Is(err, service.ErrRefReplyTo):
			errUnprocessableEntityResp(c, "replyTo is invalid chatId")
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "chat media uploaded",
		"data": gin.H{
			"chat": gin.H{
				"id": dtoChat.Id,
			},
		},
	})
}

func (chat *chat) getMessagesHandler(c *gin.Context) {
	var query struct {
		Limit  *int       `form:"limit" binding:"omitempty,min=30,max=90"`
		At     *time.Time `form:"at" binding:"omitempty,required_with_all=ChatId After"`
		ChatId *string    `form:"chatId" binding:"omitempty,required_with_all=At After,uuid"`
		After  *bool      `form:"after" binding:"omitempty,required_with_all=At ChatId"`
	}

	if err := c.ShouldBindQuery(&query); err != nil {
		if errMap := util.ReadValidationErr(err, map[string]string{
			"Limit":  "if provided, value must in between 30-90",
			"At":     "if provided, must provide chatId, after field also. must be valid time",
			"ChatId": "if provided, must provide at, after field also. must be valid uuid",
			"After":  "if provided, must provide ahatId, after field also",
		}); errMap != nil {
			errValidationResp(c, errMap)
			return
		}
		errServerResp(c, err)
		return
	}
	chatQueryFilter := entity.ChatFilter{}
	if query.Limit != nil {
		chatQueryFilter.Limit = *query.Limit
	}
	if query.At != nil {
		chatQueryFilter.Cursor = &entity.ChatCursor{
			At: *query.At,
		}
	}
	if query.ChatId != nil {
		chatQueryFilter.Cursor.Id = *query.ChatId
	}
	if query.After != nil {
		chatQueryFilter.Cursor.After = *query.After
	}
	convoId := c.GetString("convId")
	dtoChats, err := chat.chatSvc.GetMessages(convoId, chatQueryFilter)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccessingDB) {
			errResourceConflictResp(c)
			return
		}
		errServerResp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"chats": dtoChats,
		},
	})
}

func (chat *chat) deleteMessagesByIdHandler(c *gin.Context) {
	chatId := c.GetString("chatId")
	if err := chat.chatSvc.DeleteMessagesById(chatId); err != nil {
		switch {
		case errors.Is(err, domain.ErrRefNotFound23503):
			errNotFoundResp(c, "provided chatId in url is not found!")
		case errors.Is(err, domain.ErrTooLongAccessingDB):
			errResourceConflictResp(c)
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"messages": "chat deleted",
	})

}
