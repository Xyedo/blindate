package dtos

import (
	"errors"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/domain/conversation/entities"
	"github.com/xyedo/blindate/pkg/optional"
	"github.com/xyedo/blindate/pkg/pagination"
)

type IndexConverastionQueryParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (params *IndexConverastionQueryParams) Mod() *IndexConverastionQueryParams {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	return params
}

func (params IndexConverastionQueryParams) Validate() error {
	return validation.ValidateStruct(&params,
		validation.Field(&params.Page, validation.Required, validation.Min(1)),
		validation.Field(&params.Limit, validation.Required, validation.Min(1)),
	)
}

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
