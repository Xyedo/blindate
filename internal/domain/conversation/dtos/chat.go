package dtos

import (
	"errors"
	"time"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/domain/conversation/entities"
	"github.com/xyedo/blindate/pkg/optional"
	"github.com/xyedo/blindate/pkg/pagination"
)

type IndexChatQueryParams struct {
	Limit int     `query:"limit"`
	Prev  *string `query:"prev"`
	Next  *string `query:"next"`
}

func (params *IndexChatQueryParams) Mod() *IndexChatQueryParams {
	if params.Limit == 0 {
		params.Limit = 10
	}

	return params
}

func (params IndexChatQueryParams) Validate() error {
	err := validation.ValidateStruct(&params,
		validation.Field(&params.Prev, validation.NilOrNotEmpty, is.Base64, validation.By(validateCursor)),
		validation.Field(&params.Next, validation.NilOrNotEmpty, is.Base64, validation.By(validateCursor)),
		validation.Field(&params.Limit, validation.Required, validation.Min(1)),
	)

	if err != nil {
		return err
	}

	if params.Prev != nil && params.Next != nil {
		return apperror.BadPayloadWithPayloadMap(apperror.PayloadMap{
			Payloads: []apperror.ErrorPayload{
				{
					Code:    apperror.StatusErrorValidation,
					Message: "both value present",
					Details: map[string][]string{
						"prev": {"should not present when next is present"},
						"next": {"should not present when prev is present"},
					},
				},
			},
		})
	}

	return nil
}

func (params IndexChatQueryParams) ToEntity(requestId, conversationId string) entities.IndexChatPayload {
	var next, prev entities.IndexChatPayloadCursor

	if params.Next != nil {
		c, err := pagination.NewCursorFromBase64(*params.Next)
		if err != nil {
			panic("must be validate first")
		}

		next.ChatId = optional.NewString(c.Id)
		next.SentAt = optional.NewTime(c.Date)
	}

	if params.Prev != nil {
		c, err := pagination.NewCursorFromBase64(*params.Prev)
		if err != nil {
			panic("must be validate first")
		}

		prev.ChatId = optional.NewString(c.Id)
		prev.SentAt = optional.NewTime(c.Date)
	}

	return entities.IndexChatPayload{
		RequestId:      requestId,
		ConversationId: conversationId,
		Limit:          params.Limit,
		Next:           next,
		Prev:           prev,
	}
}

func validateCursor(value any) error {
	value, isNil := validation.Indirect(value)
	if isNil || validation.IsEmpty(value) {
		return nil
	}

	_, err := pagination.NewCursorFromBase64(value.(string))
	if err != nil {
		if errors.Is(err, pagination.ErrInvalidCursorFormat) {
			return errors.New("invalid format")
		}
		return err
	}

	return nil
}

type IndexChatPayload struct {
	HasNext, HasPrev bool
	Conv             entities.Conversation
}
type IndexChatConversation struct {
	MatchId   string             `json:"id"`
	ChatRows  int64              `json:"chat_rows"`
	DayPass   int64              `json:"day_pass"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Version   int64              `json:"-"`
	Recepient IndexChatUser      `json:"recepient"`
	Chats     []IndexChatElement `json:"chats"`
}

type IndexChatUser struct {
	Id          string          `json:"id"`
	DisplayName string          `json:"display_name"`
	FileId      optional.String `json:"-"`
	Url         string          `json:"url"`
}

type IndexChatElement struct {
	Id             string          `json:"id"`
	ConversationId string          `json:"conversation_id"`
	Author         string          `json:"author"`
	Messages       string          `json:"message"`
	ReplyTo        optional.String `json:"reply_to"`
	SentAt         time.Time       `json:"sent_at"`
	SeenAt         optional.Time   `json:"seen_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Version        int             `json:"-"`
}
type IndexChatMetadata struct {
	Next *string `json:"next"`
	Prev *string `json:"prev"`
}
type IndexChatResponse struct {
	Metadata IndexChatMetadata     `json:"metadata"`
	Data     IndexChatConversation `json:"data"`
}

func NewIndexChatResponse(payload IndexChatPayload) IndexChatResponse {
	chats := make([]IndexChatElement, 0, len(payload.Conv.Chats))
	for _, chat := range payload.Conv.Chats {
		chats = append(chats, IndexChatElement(chat))
	}

	var next, prev *string
	if len(chats) > 0 {
		if payload.HasNext {
			lastChat := chats[len(chats)-1]
			*next = pagination.NewBase64FromCursor(pagination.Cursor{
				Id:   lastChat.Id,
				Date: lastChat.SentAt,
			})
		}

		if payload.HasPrev {
			firstChat := chats[0]
			*prev = pagination.NewBase64FromCursor(pagination.Cursor{
				Id:   firstChat.Id,
				Date: firstChat.SentAt,
			})
		}
	}
	return IndexChatResponse{
		Metadata: IndexChatMetadata{
			Next: next,
			Prev: prev,
		},
		Data: IndexChatConversation{
			MatchId:   payload.Conv.MatchId,
			ChatRows:  payload.Conv.ChatRows,
			DayPass:   payload.Conv.DayPass,
			CreatedAt: payload.Conv.CreatedAt,
			UpdatedAt: payload.Conv.UpdatedAt,
			Version:   payload.Conv.Version,
			Recepient: IndexChatUser(payload.Conv.Recepient),
			Chats:     chats,
		},
	}
}
