package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/entity"
)

func NewChat(conn *sqlx.DB) *chat {
	return &chat{
		conn,
	}
}

type chat struct {
	*sqlx.DB
}

// TODO: create chat_test.go
func (c *chat) InsertNewChat(content *entity.Chat) error {
	chatQ := `
	INSERT INTO chats(conversation_id,messages,reply_to,sent_at)
	VALUES($1,$2,$3)
	RETURNING id`
	contentArgs := []any{content.ConversationId, content.Messages, content.ReplyTo, content.SentAt}

	attachmentQ := `
	INSERT INTO media(chat_id, blob_link,media_type)
	VALUES($1,$2,$3)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.execTx(ctx, func(q *sqlx.DB) error {
		err := c.GetContext(ctx, &content.Id, chatQ, contentArgs...)
		if err != nil {
			return err
		}
		if content.Attachment != nil {
			attachmentArgs := []any{content.Id, content.Attachment.BlobLink, content.Attachment.MediaType}
			_, err = c.ExecContext(ctx, attachmentQ, attachmentArgs...)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *chat) DeleteChat(chatId string) (int64, error) {
	query := `
	DELETE FROM chats WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.ExecContext(ctx, query, chatId)
	if err != nil {
		return 0, err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return row, nil
}

func (c *chat) SelectChat(convoId string, filter entity.ChatFilter) ([]entity.Chat, error) {
	query := `
	SELECT 
		chats.*
		media.blob_link
		media.media_type
	FROM chats
	LEFT JOIN media
		ON chats.id = media.chat_id
	WHERE 
		chats.conversation_id=$1`
	if filter.Cursor != nil {
		if filter.Cursor.After {
			query += ` (AND chats.sent_at > $2  AND chats.id = $3`
		} else {
			query += ` (AND chats.sent_at <= $2 AND chats.id = $3`
		}
		query += ` OR $2 = '' AND $3 = '')`
	}
	query += ` ORDER BY chats.sent_at DESC
	LIMIT $4`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var args []any

	if filter.Cursor == nil {
		args = []any{convoId, "", "", filter.Limit}
	} else {
		args = []any{convoId, filter.Cursor.At, filter.Cursor.Id, filter.Limit}
	}

	chats := []entity.Chat{}
	err := c.SelectContext(ctx, &chats, query, args...)
	if err != nil {
		return nil, err
	}
	return chats, nil
}
func (c *chat) execTx(ctx context.Context, q func(q *sqlx.DB) error) error {
	return execGeneric(c.DB, ctx, q, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}
