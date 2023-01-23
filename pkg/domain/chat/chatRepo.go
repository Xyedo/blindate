package chat

import (
	"time"

	chatEntity "github.com/xyedo/blindate/pkg/domain/chat/entities"
)

type Filter struct {
	Cursor *Cursor
	Limit  int
}
type Cursor struct {
	At    time.Time
	Id    string
	After bool
}
type Repository interface {
	InsertNewChat(content *chatEntity.DAO) error
	SelectChat(convoId string, filter Filter) ([]chatEntity.DAO, error)
	UpdateSeenChat(convId, authorId string) ([]string, error)
	DeleteChatById(chatId string) error
}
