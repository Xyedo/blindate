package entity

import (
	"database/sql"
	"time"

	"github.com/xyedo/blindate/pkg/domain"
)

type Chat struct {
	Id             string             `db:"id"`
	ConversationId string             `db:"conversation_id"`
	Messages       string             `db:"messages"`
	ReplyTo        sql.NullString     `db:"reply_to"`
	SentAt         time.Time          `db:"sent_at"`
	SeenAt         sql.NullTime       `db:"seen_at"`
	Attachment     *domain.Attachment `db:"attachment"`
}

type ChatFilter struct {
	Cursor *cursor
	Limit  int
}
type cursor struct {
	At    time.Time
	Id    string
	After bool
}
