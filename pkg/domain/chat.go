package domain

import (
	"time"
)

type Chat struct {
	Id             string      `json:"id"`
	ConversationId string      `json:"conversationId"`
	Messages       string      `json:"messages"`
	ReplyTo        *string     `json:"replyTo"`
	SentAt         time.Time   `json:"sentAt"`
	SeenAt         *time.Time  `json:"seenAt"`
	Attachment     *Attachment `json:"attachment"`
}

type Attachment struct {
	ChatId    string `json:"-" db:"chat_id"`
	BlobLink  string `json:"blobLink" db:"blob_link"`
	MediaType string `json:"mediaType" db:"media_type"`
}
