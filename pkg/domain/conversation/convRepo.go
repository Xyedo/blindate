package conversation

import (
	convEntity "github.com/xyedo/blindate/pkg/domain/conversation/entities"
)

type Filter struct {
	Offset int
}
type Repository interface {
	InsertConversation(matchId string) (string, error)
	SelectConversationById(matchId string) (convEntity.DTO, error)
	SelectConversationByUserId(UserId string, filter *Filter) ([]convEntity.DTO, error)
	UpdateDayPass(convoId string) error
	UpdateChatRow(convoId string) error
	DeleteConversationById(convoId string) error
}
