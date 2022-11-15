package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/domain"
)

type ConversationRepo interface {
	InsertConversation(convo *domain.Conversation) error
	// SelectConversationById(convoId string) (*domain.Conversation, error)
	SelectConversationByUserId(UserId string) ([]domain.Conversation, error)
	UpdateDayPass(convoId string) error
	UpdateChatRow(convoId string) error
}

func NewConversation(conn *sqlx.DB) *conversation {
	return &conversation{
		conn,
	}
}

type conversation struct {
	*sqlx.DB
}

// TODO: create conversation.go
func (c *chat) InsertConversation(convo *domain.Conversation) error {
	query := `
	INSERT INTO conversations(from_id,to_id)
	VALUES($1,$2)
	RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.GetContext(ctx, &convo.Id, query, convo.FromUser.ID, convo.ToUser.ID)
	if err != nil {
		return err
	}
	return nil
}

// func (c *chat) SelectConversationById(convoId string) (*domain.Conversation, error) {
// 	query := `
// 	SELECT
// 		convo.id AS id,
// 		convo.chat_rows AS chat_rows,
// 		convo.day_pass AS day_pass,
// 		creator.*,
// 		recipient.*
// 		c.messages AS last_messages,
// 		c.sent_at AS last_messages_sent_at,
// 		c.seen_at AS last_messages_seen_at
// 	FROM conversations AS convo
// 	JOIN users AS creator
// 		ON creator.id = convo.from_id
// 	JOIN users AS recipient
// 		ON recipient.id = convo.to_id
// 	LEFT JOIN (
// 		SELECT DISTINCT ON (conversation_id)
// 			messages,
// 			sent_at,
// 			seen_at
// 		FROM chats
// 		ORDER BY conversation_id, sent_at DESC
// 		) chats AS c ON c.conversation_id = convo.id
// 	WHERE convo.id = $1`

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	var conv *domain.Conversation
// 	err := c.GetContext(ctx, &conv, query, convoId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return conv, nil
// }

// LEFT JOIN (
// 	SELECT
// 		id,
// 		user_id,
// 		selected,
// 		picture_ref
// 	FROM profile_picture
// 	WHERE selected = true

// 	UNION

// 	SELECT
// 		id,
// 		user_id,
// 		selected,
// 		picture_ref
// 	FROM profile_picture
// 	ORDER BY id ASC
// 	LIMIT 1

// ) profile_picture as creator_pp
// 		ON creator_pp.user_id = creator.id
// LEFT JOIN (
// 	SELECT
// 		id,
// 		user_id,
// 		selected,
// 		picture_ref
// 	FROM profile_picture
// 	WHERE selected = true

// 	UNION

// 	SELECT
// 		id,
// 		user_id,
// 		selected,
// 		picture_ref
// 	FROM profile_picture
// 	ORDER BY id ASC
// 	LIMIT 1

// ) profile_picture as recipient_pp
//
//	ON recipient_pp.user_id = recipient.id

func (c *chat) SelectConversationByUserId(UserId string) ([]domain.Conversation, error) {
	convQuery := `
	SELECT 
		convo.id AS id,
		convo.chat_rows AS chat_rows,
		convo.day_pass AS day_pass,
		creator.*,
		recipient.*,
		c.messages AS last_messages,
		c.sent_at AS last_messages_sent_at,
		c.seen_at AS last_messages_seen_at
	FROM conversations AS conv
	LEFT JOIN (
		( 
			SELECT
				u.id as id,
				u.full_name as full_name,
				u.alias as alias,
			pp.picture_ref as picture_ref
			FROM users AS u
			LEFT JOIN profile_picture as pp
				ON pp.user_id = u.id
			WHERE selected = true
		)
		UNION ALL
		(
			SELECT
				u.id as id,
				u.full_name as full_name,
				u.alias as alias,
				pp.picture_ref as picture_ref
			FROM users AS u
			LEFT JOIN profile_picture as pp
				ON pp.user_id = u.id
			ORDER BY id ASC
			LIMIT 1
		)
		LIMIT 1
	) users AS creator
		ON creator.id = conv.from_id
	LEFT JOIN (
		(
			SELECT
				u.id as id,
				u.full_name as full_name,
				u.alias as alias,
				pp.picture_ref as picture_ref
			FROM users AS u
			LEFT JOIN profile_picture as pp
				ON pp.user_id = u.id
			WHERE selected = true
		)
		UNION ALL
		(
			SELECT
				u.id as id,
				u.full_name as full_name,
				u.alias as alias,
				pp.picture_ref as picture_ref
			FROM users AS u
			LEFT JOIN profile_picture as pp
				ON pp.user_id = u.id
			ORDER BY id ASC
			LIMIT 1
		)
		LIMIT 1
	) users AS recipient
		ON recipient.id = conv.to_id
	LEFT JOIN (
		SELECT DISTINCT ON (conversation_id) 
			messages,
			sent_at,
			seen_at
		FROM chats 
		ORDER BY conversation_id, sent_at DESC
		) chats AS c ON c.conversation_id = convo.id
	WHERE 
		creator.id = $1 OR
		recipient.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	convs := []domain.Conversation{}
	err := c.SelectContext(ctx, &convs, convQuery, UserId)
	if err != nil {
		return nil, err
	}
	return convs, nil
}
func (c *chat) UpdateChatRow(convoId string) error {
	query := `
	UPDATE conversations SET
		chat_rows = chat_rows + 1
	WHERE id = $1
	RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := c.GetContext(ctx, &id, query, convoId)
	if err != nil {
		return err
	}
	return nil
}

func (c *chat) UpdateDayPass(convoId string) error {
	query := `
	UPDATE conversations SET
		day_pass = day_pass +1
	WHERE id = $1
	RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id string
	err := c.GetContext(ctx, &id, query, convoId)
	if err != nil {
		return err
	}
	return nil
}
