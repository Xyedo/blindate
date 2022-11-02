package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/domain"
)

func NewConversation(conn *sqlx.DB) *conversation {
	return &conversation{
		conn,
	}
}

type conversation struct {
	*sqlx.DB
}

func (c *conversation) InsertConversation(convo *domain.Conversation) error {
	query := `
	INSERT INTO conversations(from_id,to_id)
	VALUES($1,$2)
	RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.GetContext(ctx, &convo.Id, query, convo.FromId, convo.ToId)
	if err != nil {
		return err
	}
	return nil
}

func (c *conversation) SelectConversation(convoId string) (*domain.Conversation, error) {
	//TODO: left Join message one->get lastMessage, sent_at, or seen_at
	query := `
	SELECT 
		id,
		from_id,
		to_id,
		chat_rows,
		day_pass
	FROM conversations
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var conv *domain.Conversation
	err := c.GetContext(ctx, &conv, query, convoId)
	if err != nil {
		return nil, err
	}
	return conv, nil
}
func (c *conversation) UpdateChatRow(convoId string) error {
	query := `
	UPDATE conversations SET
		chat_rows = chat_rows + 1
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.ExecContext(ctx, query, convoId)
	if err != nil {
		return err
	}
	return nil
}

func (c *conversation) UpdateDayPass(convoId string) error {
	query := `
	UPDATE conversations SET
		day_pass = day_pass +1
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := c.ExecContext(ctx, query, convoId)
	if err != nil {
		return err
	}
	return nil
}
