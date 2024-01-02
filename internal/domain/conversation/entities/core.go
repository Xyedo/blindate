package entities

import (
	"time"

	"github.com/xyedo/blindate/pkg/optional"
)

type Conversation struct {
	MatchId   string
	ChatRows  int64
	DayPass   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64

	//computed value
	Recepient RecepientUser
	Chats     []Chat
}

type Chat struct {
	Id             string
	ConversationId string
	Author         string
	Messages       string
	ReplyTo        optional.String
	SentAt         time.Time
	SeenAt         optional.Time
	UpdatedAt      time.Time
	Version        int
}
