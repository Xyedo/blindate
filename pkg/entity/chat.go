package entity

import (
	"database/sql"
	"time"
)

type Chat struct {
	Id             string         `db:"id"`
	ConversationId string         `db:"conversation_id"`
	Messages       string         `db:"messages"`
	ReplyTo        sql.NullString `db:"reply_to"`
	SentAt         time.Time      `db:"sent_at"`
	SeenAt         sql.NullTime   `db:"seen_at"`
	Attachment     *Attachment    `db:"attachment"`
}

type Attachment struct {
	ChatId    string `db:"chat_id"`
	BlobLink  string `db:"blob_link"`
	MediaType string `db:"media_type"`
}
