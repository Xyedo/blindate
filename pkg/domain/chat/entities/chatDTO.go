package chatEntity

import "time"

// Chat one convoersation to many chat
type DTO struct {
	Id             string      `json:"id"`
	ConversationId string      `json:"conversationId"`
	Author         string      `json:"author"`
	Messages       string      `json:"messages"`
	ReplyTo        *string     `json:"replyTo"`
	SentAt         time.Time   `json:"sentAt"`
	SeenAt         *time.Time  `json:"seenAt"`
	Attachment     *Attachment `json:"attachment"`
}

// Attachment one to one with chat
type Attachment struct {
	ChatId    string `json:"-" db:"chat_id"`
	BlobLink  string `json:"blobLink" db:"blob_link"`
	MediaType string `json:"mediaType" db:"media_type"`
}
