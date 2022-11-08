package domain

import (
	"time"
)

type Chat struct {
	Id             string     `json:"id"`
	ConversationId string     `json:"conversationId"`
	Messages       string     `json:"messages"`
	ReplyTo        *string    `json:"replyTo"`
	SentAt         time.Time  `json:"sentAt"`
	SeenAt         *time.Time `json:"seenAt"`
}
