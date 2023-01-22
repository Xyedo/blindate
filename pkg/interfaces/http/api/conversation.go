package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

// TODO: Create conversation_test.go
type conversationSvc interface {
	CreateConversation(matchId string) (string, error)
	FindConversationById(convoId string) (domain.Conversation, error)
	GetConversationByUserId(userId string) ([]domain.Conversation, error)
	DeleteConversationById(convoId string) error
}

func NewConvo(convSvc conversationSvc) *Conversation {
	return &Conversation{
		convSvc: convSvc,
	}
}

type Conversation struct {
	convSvc conversationSvc
}

func (conv *Conversation) postConversationHandler(c *gin.Context) {
	var input struct {
		MatchId string `json:"matchId" binding:"required,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"matchId": "must be required and must valid uuid",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	_, err = conv.convSvc.CreateConversation(input.MatchId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"conversationId": input.MatchId,
		},
	})
}

func (conv *Conversation) getConversationByUserId(c *gin.Context) {
	userId := c.GetString(keyUserId)
	convs, err := conv.convSvc.GetConversationByUserId(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"conversations": convs,
		},
	})
}
func (conv *Conversation) getConversationById(c *gin.Context) {
	convId := c.GetString(keyConvId)
	convRet, err := conv.convSvc.FindConversationById(convId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"conversation": convRet,
		},
	})
}
func (conv *Conversation) deleteConversationById(c *gin.Context) {
	convId := c.GetString(keyConvId)
	if err := conv.convSvc.DeleteConversationById(convId); err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "deleting conversationId success",
	})
}
