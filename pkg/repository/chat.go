package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
)

type ChatRepo interface {
	InsertNewChat(content *entity.Chat) error
	SelectChat(convoId string, filter entity.ChatFilter) ([]entity.Chat, error)
	UpdateSeenChatById(chatId string) error
	DeleteChatById(chatId string) error
}

func NewChat(conn *sqlx.DB) *chat {
	return &chat{
		conn: conn,
	}
}

type chat struct {
	conn *sqlx.DB
}

func (c *chat) InsertNewChat(content *entity.Chat) error {
	chatQ := `
	INSERT INTO chats(conversation_id,author,messages,reply_to,sent_at)
	VALUES($1,$2,$3,$4, $5)
	RETURNING id`
	contentArgs := []any{
		content.ConversationId,
		content.Author,
		content.Messages,
		content.ReplyTo,
		content.SentAt,
	}

	attachmentQ := `
	INSERT INTO media(chat_id, blob_link,media_type)
	VALUES($1,$2,$3)`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.execTx(ctx, func(q *sqlx.DB) error {
		err := q.GetContext(ctx, &content.Id, chatQ, contentArgs...)
		if err != nil {
			return err
		}
		if content.Attachment != nil {
			attachmentArgs := []any{content.Id, content.Attachment.BlobLink, content.Attachment.MediaType}
			_, err = q.ExecContext(ctx, attachmentQ, attachmentArgs...)
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

func (c *chat) DeleteChatById(chatId string) error {
	query := `
	DELETE FROM chats WHERE id = $1 RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var retChatId string
	err := c.conn.GetContext(ctx, &retChatId, query, chatId)
	if err != nil {
		return err
	}

	return nil
}

// func (c *chat) SelectChatById(chatId string) (*entity.Chat, error) {
// 	query := `
// 	SELECT
// 		chats.id,
// 		chats.conversation_id,
// 		chats.author,
// 		chats.messages,
// 		chats.reply_to,
// 		chats.sent_at,
// 		chats.seen_at,
// 		media.blob_link,
// 		media.media_type
// 	FROM chats
// 	LEFT JOIN media
// 		ON media.chat_id = chats.id
// 	WHERE chats.id = $1`
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

//		row := c.conn.QueryRowxContext(ctx, query, chatId)
//		newChat, err := c.createNewChat(row)
//		if err != nil {
//			return nil, err
//		}
//		return &newChat, nil
//	}
func (c *chat) UpdateSeenChatById(chatId string) error {
	query := `UPDATE chats SET seen_at = $1 WHERE id = $2 RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var retId string

	err := c.conn.GetContext(ctx, &retId, query, sql.NullTime{Valid: true, Time: time.Now()}, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (c *chat) SelectChat(convoId string, filter entity.ChatFilter) ([]entity.Chat, error) {
	if filter.Limit == 0 {
		filter.Limit = 20
	}
	args := make([]any, 0)
	query := `
	SELECT 
		chats.id,
		chats.conversation_id,
		chats.author,
		chats.messages,
		chats.reply_to,
		chats.sent_at,
		chats.seen_at,
		media.blob_link,
		media.media_type
	FROM chats
	LEFT JOIN media
		ON chats.id = media.chat_id
	WHERE 
		chats.conversation_id=$1`
	args = append(args, convoId)
	if filter.Cursor != nil {
		if filter.Cursor.After {
			query += ` AND (chats.id, chats.sent_at) < ($2, $3::TIMESTAMPTZ)`
		} else {
			query += ` AND (chats.id, chats.sent_at) > ($2, $3::TIMESTAMPTZ)`
		}
		args = append(args, filter.Cursor.Id, filter.Cursor.At)
		query += ` ORDER BY chats.sent_at DESC, id DESC
		LIMIT $4`
	} else {
		query += ` ORDER BY chats.sent_at DESC, id DESC
		LIMIT $2`
	}
	args = append(args, filter.Limit)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := c.conn.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	chats := make([]entity.Chat, 0)
	for rows.Next() {
		newChat, err := c.createNewChat(rows)
		if err != nil {
			return nil, err
		}
		chats = append(chats, newChat)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return chats, nil
}
func (*chat) createNewChat(row sqlx.ColScanner) (entity.Chat, error) {
	var newChat entity.Chat
	var blobLink sql.NullString
	var mediaType sql.NullString
	err := row.Scan(
		&newChat.Id,
		&newChat.ConversationId,
		&newChat.Author,
		&newChat.Messages,
		&newChat.ReplyTo,
		&newChat.SentAt,
		&newChat.SeenAt,
		&blobLink,
		&mediaType,
	)
	if err != nil {
		return entity.Chat{}, err
	}
	if blobLink.Valid && mediaType.Valid {
		newChat.Attachment = &domain.ChatAttachment{
			ChatId:    newChat.Id,
			BlobLink:  blobLink.String,
			MediaType: mediaType.String,
		}
	}
	return newChat, nil
}
func (c *chat) execTx(ctx context.Context, q func(q *sqlx.DB) error) error {
	return execGeneric(c.conn, ctx, q, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}
