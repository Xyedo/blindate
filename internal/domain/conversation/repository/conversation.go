package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/domain/conversation/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
	"github.com/xyedo/blindate/pkg/pagination"
)

func CreateConversation(ctx context.Context, conn pg.Querier, payload entities.Conversation) error {
	const createConversation = `
		INSERT INTO conversations (
			match_id,
			chat_rows,
			day_pass,
			created_at,
			updated_at,
			version
		)
		VALUES($1,$2,$3,$4,$5,$6)
		returning match_id
	`
	var returningMatchId string
	err := conn.
		QueryRow(ctx, createConversation,
			payload.MatchId,
			payload.ChatRows,
			payload.DayPass,
			payload.CreatedAt,
			payload.UpdatedAt,
			payload.Version,
		).Scan(&returningMatchId)
	if err != nil {
		return err
	}

	if returningMatchId != payload.MatchId {
		return errors.New("som wen wong")
	}

	return nil
}

// FindConversationsByUserId
// return conversationIndex, hasNext for pagination, error
func FindConversationsByUserId(ctx context.Context, conn pg.Querier, userId string, pagination pagination.Pagination) (entities.ConversationIndex, bool, error) {
	const findConversationsByUserId = `
	SELECT
		conv.match_id,
		CASE WHEN m.request_from = $1 THEN m.request_to ELSE m.request_from END,
		ad.alias,
		pp.file_id,
		conv.chat_rows,
		conv.day_pass,
		conv.created_at,
		conv.updated_at,
		conv.version,
		last_chat.id,
		last_chat.author,
		last_chat.messages,
		last_chat.unread_message_count,
		last_chat.reply_to,
		last_chat.sent_at,
		last_chat.seen_at,
		last_chat.updated_at,
		last_chat.version
	FROM conversations conv
	JOIN match m ON conv.match_id = m.id
	JOIN account_detail ad ON ad.account_id = CASE WHEN m.request_from = $1 THEN m.request_to ELSE m.request_from END
	LEFT JOIN profile_pictures pp ON pp.account_id = ad.account_id AND pp.selected = TRUE
	LEFT JOIN (
		SELECT DISTINCT on (c1.id) 
			c1.id,
			c1.conversation_id,
			c1.author,
			c1.messages,
			c1.reply_to,
			c1.sent_at,
			c1.seen_at,
			c1.updated_at,
			c1.version,
			c2.unread_message_count
		from chat c1 JOIN (
			select 
				conversation_id,
				COUNT(CASE WHEN seen_at IS NULL THEN 1 END) as unread_message_count
			FROM chat
			GROUP BY conversation_id
		)  c2 on c2.conversation_id = c1.conversation_id
		ORDER BY id, sent_at desc
	) as last_chat ON conv.match_id = last_chat.conversation_id
	WHERE 
		m.request_from = $1 OR m.request_to = $1
	ORDER BY last_chat.sent_at desc
	OFFSET $2
	LIMIT $3
	`

	rows, err := conn.Query(ctx, findConversationsByUserId, userId, pagination.Offset(), pagination.Limit+1)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	conversationIndex := make([]entities.ConversationElement, 0)
	for rows.Next() {
		var conversation entities.ConversationElement
		err = rows.Scan(
			&conversation.Id,
			&conversation.Recepient.Id,
			&conversation.Recepient.DisplayName,
			&conversation.Recepient.FileId,
			&conversation.ChatRows,
			&conversation.DayPass,
			&conversation.CreatedAt,
			&conversation.UpdatedAt,
			&conversation.Version,
			&conversation.LastChat.Id,
			&conversation.LastChat.Author,
			&conversation.LastChat.Message,
			&conversation.LastChat.UnreadMessageCount,
			&conversation.LastChat.ReplyTo,
			&conversation.LastChat.SentAt,
			&conversation.LastChat.SeenAt,
			&conversation.LastChat.UpdatedAt,
			&conversation.LastChat.Version,
		)
		if err != nil {
			return nil, false, err
		}

		conversationIndex = append(conversationIndex, conversation)
	}

	var hasNext bool
	if len(conversationIndex) > pagination.Limit {
		hasNext = true
	}

	return conversationIndex[:pagination.Limit], hasNext, nil
}

// FindConversationsByUserId
// return conversation, hasNext, hasPrev for pagination, error
func FindChatsByConversationId(ctx context.Context, conn pg.Querier, payload entities.IndexChatPayload) (entities.Conversation, bool, bool, error) {
	const findConverastionById = `
	SELECT
		conv.match_id,
		CASE WHEN m.request_from = $1 THEN m.request_to ELSE m.request_from END,
		ad.alias,
		pp.file_id,
		conv.chat_rows,
		conv.day_pass,
		conv.created_at,
		conv.updated_at,
		conv.version
	FROM conversations conv
	JOIN match m ON conv.match_id = m.id
	JOIN account_detail ad ON ad.account_id = CASE WHEN m.request_from = $1 THEN m.request_to ELSE m.request_from END
	LEFT JOIN profile_pictures pp ON pp.account_id = ad.account_id AND pp.selected = TRUE
	WHERE 
		conv.match_id = $2
	`
	var batch pgx.Batch
	var conversation entities.Conversation
	batch.
		Queue(findConverastionById, payload.RequestId, payload.ConversationId).
		QueryRow(func(row pgx.Row) error {
			err := row.Scan(
				&conversation.MatchId,
				&conversation.Recepient.Id,
				&conversation.Recepient.DisplayName,
				&conversation.Recepient.FileId,
				&conversation.ChatRows,
				&conversation.DayPass,
				&conversation.CreatedAt,
				&conversation.UpdatedAt,
				&conversation.Version,
			)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return nil
				}
				return err
			}

			return nil
		})

	const findChatsByConversationId = `
		SELECT 
			id,
			conversation_id,
			author,
			messages,
			reply_to,
			sent_at,
			seen_at,
			updated_at,
			version
		FROM chat
		WHERE 
			conversation_id = $1 
			AND (sent_at, id) < ($3,$2)
			AND (sent_at, id) > ($5,$4)
		LIMIT $6
		ORDER BY sent_at DESC, id DESC
	`
	batch.
		Queue(
			findChatsByConversationId,
			payload.ConversationId,
			payload.Next.ChatId,
			payload.Next.SentAt,
			payload.Prev.ChatId,
			payload.Prev.SentAt,
			payload.Limit,
		).
		Query(func(rows pgx.Rows) error {
			chats := make([]entities.Chat, 0)
			for rows.Next() {
				var chat entities.Chat
				err := rows.Scan(
					&chat.Id,
					&chat.ConversationId,
					&chat.Author,
					&chat.Messages,
					&chat.ReplyTo,
					&chat.SentAt,
					&chat.SeenAt,
					&chat.UpdatedAt,
					&chat.Version,
				)
				if err != nil {
					return err
				}

				chats = append(chats, chat)
			}

			conversation.Chats = chats

			return nil
		})

	var hasNext bool
	if payload.Next.ChatId.IsPresent() && payload.Next.SentAt.IsPresent() {
		const checkHasNext = `
			SELECT count(1)
		FROM chat
		WHERE 
			conversation_id = $1 
			AND (sent_at, id) < ($3,$2)
		LIMIT $4
		ORDER BY sent_at DESC, id DESC
		`

		batch.
			Queue(checkHasNext,
				payload.ConversationId,
				payload.Next.ChatId,
				payload.Next.SentAt,
				payload.Limit,
			).
			QueryRow(func(row pgx.Row) error {
				var count int64
				err := row.Scan(&count)
				if err != nil {
					return err
				}

				hasNext = count != 0
				return nil
			})

	}

	var hasPrev bool
	if payload.Prev.ChatId.IsPresent() && payload.Prev.SentAt.IsPresent() {
		const checkHasPrev = `
			SELECT count(1)
		FROM chat
		WHERE 
			conversation_id = $1 
			AND (sent_at, id) > ($3,$2)
		LIMIT $4
		ORDER BY sent_at DESC, id DESC
		`

		batch.
			Queue(checkHasPrev,
				payload.ConversationId,
				payload.Prev.ChatId,
				payload.Prev.SentAt,
				payload.Limit,
			).
			QueryRow(func(row pgx.Row) error {
				var count int64
				err := row.Scan(&count)
				if err != nil {
					return err
				}

				hasPrev = count != 0
				return nil
			})

	}

	err := conn.SendBatch(ctx, &batch).Close()
	if err != nil {
		return entities.Conversation{}, false, false, err
	}

	return conversation, hasNext, hasPrev, nil

}
