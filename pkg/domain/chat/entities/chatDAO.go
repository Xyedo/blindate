package chatEntity

import (
	"database/sql"
	"time"
)

type DAO struct {
	Id             string         `db:"id"`
	ConversationId string         `db:"conversation_id"`
	Author         string         `db:"author"`
	Messages       string         `db:"messages"`
	ReplyTo        sql.NullString `db:"reply_to"`
	SentAt         time.Time      `db:"sent_at"`
	SeenAt         sql.NullTime   `db:"seen_at"`
	Attachment     *Attachment    `db:"attachment"`
}
