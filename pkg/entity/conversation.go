package entity

import (
	"database/sql"
	"time"

	"github.com/xyedo/blindate/pkg/domain"
)

type Chat struct {
	Id             string                 `db:"id"`
	ConversationId string                 `db:"conversation_id"`
	Author         string                 `db:"author"`
	Messages       string                 `db:"messages"`
	ReplyTo        sql.NullString         `db:"reply_to"`
	SentAt         time.Time              `db:"sent_at"`
	SeenAt         sql.NullTime           `db:"seen_at"`
	Attachment     *domain.ChatAttachment `db:"attachment"`
}

type ChatFilter struct {
	Cursor *ChatCursor
	Limit  int
}
type ChatCursor struct {
	At    time.Time
	Id    string
	After bool
}

type ConvFilter struct {
	Offset int
}
