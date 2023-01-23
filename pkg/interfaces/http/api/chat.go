package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	apiError "github.com/xyedo/blindate/pkg/common/error"
	"github.com/xyedo/blindate/pkg/common/util"
	"github.com/xyedo/blindate/pkg/domain/chat"
	chatEntity "github.com/xyedo/blindate/pkg/domain/chat/entities"
)

// TODO: Create chat_test.go
type chatSvc interface {
	CreateNewChat(content *chatEntity.DTO) error
	UpdateSeenChat(convId, userId string) error
	GetMessages(convoId string, filter chat.Filter) ([]chatEntity.DTO, error)
	DeleteMessagesById(chatId string) error
}

func NewChat(chatSvc chatSvc, attachSvc attachmentManager) *Chat {
	return &Chat{
		chatSvc:   chatSvc,
		attachSvc: attachSvc,
	}
}

type Chat struct {
	chatSvc   chatSvc
	attachSvc attachmentManager
}

func (cha *Chat) postChatHandler(c *gin.Context) {
	var newChatPayload chatEntity.New
	if err := c.ShouldBindJSON(&newChatPayload); err != nil {
		if jsonErr := jsonBindingErrResp(err, c, map[string]string{
			"Message": "must not empty and max characters is 4096",
			"ReplyTo": "if specified, must be valid uuid",
		}); jsonErr != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	convoId := c.GetString(keyConvId)
	userId := c.GetString("userId")
	dtoChat := chatEntity.DTO{
		ConversationId: convoId,
		Author:         userId,
		Messages:       newChatPayload.Message,
		ReplyTo:        newChatPayload.ReplyTo,
		SentAt:         time.Now(),
	}
	if err := cha.chatSvc.CreateNewChat(&dtoChat); err != nil {
		jsonHandleError(c, err)
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

func (cha *Chat) postChatMediaHandler(c *gin.Context) {
	var validAudioTypes = []string{
		"application/ogg",
		"audio/mpeg",
	}
	key, mediaType := uploadFile(c, cha.attachSvc, validAudioTypes, "chat-attachment")
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
	dtoChat := chatEntity.DTO{
		ConversationId: convoId,
		Author:         userId,
		Messages:       "",
		ReplyTo:        input.ReplyTo,
		SentAt:         time.Now(),
		Attachment: &chatEntity.Attachment{
			BlobLink:  key,
			MediaType: mediaType,
		},
	}
	if err := cha.chatSvc.CreateNewChat(&dtoChat); err != nil {
		jsonHandleError(c, err)
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

func (cha *Chat) getMessagesHandler(c *gin.Context) {
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
	chatQueryFilter := chat.Filter{}
	if query.Limit != nil {
		chatQueryFilter.Limit = *query.Limit
	}
	if query.At != nil {
		chatQueryFilter.Cursor = &chat.Cursor{
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
	dtoChats, err := cha.chatSvc.GetMessages(convoId, chatQueryFilter)
	if err != nil {
		if errors.Is(err, apiError.ErrTooLongAccessingDB) {
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
func (chat *Chat) putSeenAtHandler(c *gin.Context) {
	convoId := c.GetString("convId")
	userId := c.GetString("userId")
	err := chat.chatSvc.UpdateSeenChat(convoId, userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "seenAt updated",
	})
}
func (chat *Chat) deleteMessagesByIdHandler(c *gin.Context) {
	chatId := c.GetString("chatId")
	if err := chat.chatSvc.DeleteMessagesById(chatId); err != nil {
		switch {
		case errors.Is(err, apiError.ErrRefNotFound23503):
			errNotFoundResp(c, "provided chatId in url is not found!")
		case errors.Is(err, apiError.ErrTooLongAccessingDB):
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
