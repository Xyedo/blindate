package entities

import (
	"github.com/xyedo/blindate/pkg/optional"
)

type IndexChatPayload struct {
	RequestId      string
	ConversationId string
	Limit          int
	Next           IndexChatPayloadCursor
	Prev           IndexChatPayloadCursor
}

type IndexChatPayloadCursor struct {
	ChatId optional.String
	SentAt optional.Time
}
