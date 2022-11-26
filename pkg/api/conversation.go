package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

type conversationSvc interface {
	CreateConversation(matchId string) (string, error)
	FindConversationById(convoId string) (*domain.Conversation, error)
	GetConversationByUserId(userId string) ([]domain.Conversation, error)
	DeleteConversationById(convoId string) error
}

func NewConvo(convSvc conversationSvc) *conversation {
	return &conversation{
		convSvc: convSvc,
	}
}

type conversation struct {
	convSvc conversationSvc
}

func (conv *conversation) postConversationHandler(c *gin.Context) {
	convoId := c.GetString("convId")
	_, err := conv.convSvc.CreateConversation(convoId)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			errResourceConflictResp(c)
		case errors.Is(err, domain.ErrRefNotFound23503):
			errNotFoundResp(c, "provided userId  doesnt exists")
		case errors.Is(err, domain.ErrUniqueConstraint23505):
			errUnprocessableEntityResp(c, "already having a conversation between both userId")
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"conversationId": convoId,
		},
	})
}

func (conv *conversation) getConversationByUserId(c *gin.Context) {
	userId := c.GetString("userId")
	convs, err := conv.convSvc.GetConversationByUserId(userId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrResourceNotFound):
			errNotFoundResp(c, "conversation with this userId is not found")
		case errors.Is(err, context.Canceled):
			errResourceConflictResp(c)
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"conversations": convs,
		},
	})
}
func (conv *conversation) getConversationById(c *gin.Context) {
	convId := c.GetString("convId")
	convRet, err := conv.convSvc.FindConversationById(convId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrResourceNotFound):
			errNotFoundResp(c, "conversation with this userId is not found")
		case errors.Is(err, context.Canceled):
			errResourceConflictResp(c)
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"conversation": convRet,
		},
	})
}
func (conv *conversation) deleteConversationById(c *gin.Context) {
	convId := c.GetString("convId")
	if err := conv.convSvc.DeleteConversationById(convId); err != nil {
		switch {
		case errors.Is(err, domain.ErrResourceNotFound):
			errNotFoundResp(c, "conversationId is not found")

		case errors.Is(err, domain.ErrTooLongAccessingDB):
			errResourceConflictResp(c)
		default:
			errServerResp(c, err)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "deleting conversationId success",
	})
}
