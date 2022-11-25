package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
)

type Conversation interface {
	InsertConversation(matchId string) (string, error)
	SelectConversationById(matchId string) (*domain.Conversation, error)
	SelectConversationByUserId(UserId string, filter *entity.ConvFilter) ([]domain.Conversation, error)
	UpdateDayPass(convoId string) error
	UpdateChatRow(convoId string) error
	DeleteConversationById(convoId string) error
}

func NewConversation(conn *sqlx.DB) *conversation {
	return &conversation{
		conn: conn,
	}
}

type conversation struct {
	conn *sqlx.DB
}

func (c *conversation) InsertConversation(matchId string) (string, error) {
	query := `
	INSERT INTO conversations(match_id)
	VALUES($1)
	RETURNING match_id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var convoId string
	err := c.conn.GetContext(ctx, &convoId, query, matchId)
	if err != nil {
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
recipient.alias AS recipieint_alias,
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
JOIN match AS match
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

func (c *conversation) SelectConversationById(matchId string) (*domain.Conversation, error) {
	convQuery := selectConvo +
		` WHERE conv.match_id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := c.conn.QueryRowxContext(ctx, convQuery, matchId)
	newConv, err := c.createNewChat(row)
	if err != nil {
		return nil, err
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return newConv, nil

}

func (c *conversation) SelectConversationByUserId(UserId string, filter *entity.ConvFilter) ([]domain.Conversation, error) {
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
		return nil, err
	}
	defer rows.Close()

	convs := make([]domain.Conversation, 0)
	for rows.Next() {
		newConv, err := c.createNewChat(rows)
		if err != nil {
			return nil, err
		}
		convs = append(convs, *newConv)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return convs, nil
}
func (c *conversation) UpdateChatRow(convoId string) error {
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
		return err
	}
	return nil
}

func (c *conversation) UpdateDayPass(convoId string) error {
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
		return err
	}
	return nil
}
func (c *conversation) DeleteConversationById(convoId string) error {
	query := `
	DELETE from conversations WHERE match_id = $1 RETURNING match_id`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := c.conn.GetContext(ctx, &id, query, convoId)
	if err != nil {
		return err
	}
	return nil
}

func (c *conversation) createNewChat(row sqlx.ColScanner) (*domain.Conversation, error) {
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
		return nil, err
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
	return &newConv, nil
}
