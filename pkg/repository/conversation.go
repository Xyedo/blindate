package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/domain/entity"
)

type Conversation interface {
	InsertConversation(matchId string) (string, error)
	SelectConversationById(matchId string) (domain.Conversation, error)
	SelectConversationByUserId(UserId string, filter *entity.ConvFilter) ([]domain.Conversation, error)
	UpdateDayPass(convoId string) error
	UpdateChatRow(convoId string) error
	DeleteConversationById(convoId string) error
}

func NewConversation(conn *sqlx.DB) *ConvConn {
	return &ConvConn{
		conn: conn,
	}
}

type ConvConn struct {
	conn *sqlx.DB
}

func (c *ConvConn) InsertConversation(matchId string) (string, error) {
	query := `
	INSERT INTO conversations(match_id)
	VALUES($1)
	RETURNING match_id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var convoId string
	err := c.conn.GetContext(ctx, &convoId, query, matchId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return "", common.WrapError(err, common.ErrRefNotFound23503)
			}
			if pqErr.Code == "23505" {
				return "", common.WrapError(err, common.ErrUniqueConstraint23505)
			}
			return "", pqErr
		}
		return "", err
	}
	return convoId, nil
}

var selectConvo = `
SELECT 
conv.match_id AS id,
conv.chat_rows AS chat_rows,
conv.day_pass AS day_pass,
creator.id AS creator_id,
creator.full_name AS creator_full_name,
creator.alias AS creator_alias,
(
	SELECT
		picture_ref
	FROM profile_picture 
	WHERE user_id = creator.id
	ORDER BY selected DESC, id DESC
	LIMIT 1
) AS creator_pp_ref,
recipient.id AS recipient_id,
recipient.full_name AS recipient_full_name,
recipient.alias AS recipient_alias,
(
	SELECT 
		picture_ref
	FROM profile_picture
	WHERE user_id = recipient.id
	ORDER BY selected DESC, id DESC
	LIMIT 1
) AS recipient_pp_ref,
c.messages AS last_messages,
c.sent_at AS last_messages_sent_at,
c.seen_at AS last_messages_seen_at,
match.request_status AS request_status,
match.reveal_status AS reveal_status
FROM conversations AS conv
JOIN match
	ON match.id = conv.match_id
JOIN users AS creator 
	ON creator.id = match.request_from
JOIN users AS recipient
	ON recipient.id = match.request_to
LEFT JOIN (
	SELECT DISTINCT ON (conversation_id) 
		conversation_id,
		messages,
		sent_at,
		seen_at
	FROM chats 
	ORDER BY conversation_id, sent_at DESC
) AS c ON c.conversation_id = conv.match_id`

func (c *ConvConn) SelectConversationById(matchId string) (domain.Conversation, error) {
	convQuery := selectConvo +
		` WHERE conv.match_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := c.conn.QueryRowxContext(ctx, convQuery, matchId)
	newConv, err := c.createNewChat(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Conversation{}, common.WrapError(err, common.ErrResourceNotFound)
		}
		if errors.Is(err, context.Canceled) {
			return domain.Conversation{}, common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return domain.Conversation{}, err
	}
	if err = row.Err(); err != nil {
		return domain.Conversation{}, err
	}
	return newConv, nil

}

func (c *ConvConn) SelectConversationByUserId(UserId string, filter *entity.ConvFilter) ([]domain.Conversation, error) {
	convQuery := selectConvo +
		` WHERE 
			creator.id = $1 OR
			recipient.id = $1
		ORDER BY last_messages_sent_at DESC
		LIMIT 20`

	args := []any{UserId}
	if filter != nil {
		convQuery += ` OFFSET $2`
		args = append(args, filter.Offset)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := c.conn.QueryxContext(ctx, convQuery, args...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Panic(err)
		}
	}(rows)

	convs := make([]domain.Conversation, 0)
	for rows.Next() {
		newConv, err := c.createNewChat(rows)
		if err != nil {
			return nil, err
		}
		convs = append(convs, newConv)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return convs, nil
}
func (c *ConvConn) UpdateChatRow(convoId string) error {
	query := `
	UPDATE conversations SET
		chat_rows = chat_rows + 1
	WHERE match_id = $1
	RETURNING match_id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := c.conn.GetContext(ctx, &id, query, convoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return common.WrapError(err, common.ErrRefNotFound23503)
		}
		if errors.Is(err, context.Canceled) {
			return common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return err
	}
	return nil
}

func (c *ConvConn) UpdateDayPass(convoId string) error {
	query := `
	UPDATE conversations SET
		day_pass = day_pass +1
	WHERE match_id = $1
	RETURNING match_id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := c.conn.GetContext(ctx, &id, query, convoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return common.WrapError(err, common.ErrRefNotFound23503)
		}
		if errors.Is(err, context.Canceled) {
			return common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return err
	}
	return nil
}
func (c *ConvConn) DeleteConversationById(convoId string) error {
	query := `
	DELETE from conversations WHERE match_id = $1 RETURNING match_id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := c.conn.GetContext(ctx, &id, query, convoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return common.WrapError(err, common.ErrRefNotFound23503)
		}
		if errors.Is(err, context.Canceled) {
			return common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		return err
	}
	return nil
}

func (*ConvConn) createNewChat(row sqlx.ColScanner) (domain.Conversation, error) {
	var newConv domain.Conversation

	var creatorProfPic sql.NullString
	var recipientProfPic sql.NullString

	var lastMessage sql.NullString
	var lastMessageSentAt sql.NullTime
	var seenAt sql.NullTime
	err := row.Scan(
		&newConv.Id,
		&newConv.ChatRows,
		&newConv.DayPass,
		&newConv.FromUser.ID,
		&newConv.FromUser.FullName,
		&newConv.FromUser.Alias,
		&creatorProfPic,
		&newConv.ToUser.ID,
		&newConv.ToUser.FullName,
		&newConv.ToUser.Alias,
		&recipientProfPic,
		&lastMessage,
		&lastMessageSentAt,
		&seenAt,
		&newConv.RequestStatus,
		&newConv.RevealStatus,
	)
	if err != nil {
		return domain.Conversation{}, err
	}
	if creatorProfPic.Valid {
		newConv.FromUser.ProfilePic = creatorProfPic.String
	}
	if recipientProfPic.Valid {
		newConv.ToUser.ProfilePic = recipientProfPic.String
	}
	if lastMessage.Valid {
		newConv.LastMessage = lastMessage.String
	}
	if lastMessageSentAt.Valid {
		newConv.LastMessageSentAt = lastMessageSentAt.Time
	}
	if seenAt.Valid {
		newConv.LastMessageSeenAt = &seenAt.Time
	}
	return newConv, nil
}
