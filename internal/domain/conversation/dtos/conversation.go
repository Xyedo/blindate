package dtos

import (
	"github.com/invopop/validation"
	"github.com/xyedo/blindate/internal/domain/conversation/entities"
	"github.com/xyedo/blindate/pkg/pagination"
)

type IndexConversationQueryParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (params *IndexConversationQueryParams) Mod() *IndexConversationQueryParams {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	return params
}

func (params IndexConversationQueryParams) Validate() error {
	return validation.ValidateStruct(&params,
		validation.Field(&params.Page, validation.Required, validation.Min(1)),
		validation.Field(&params.Limit, validation.Required, validation.Min(1)),
	)
}

type IndexConversationResponse struct {
	Metadata IndexConversationMetadata      `json:"metadata"`
	Data     []entities.ConversationElement `json:"data"`
}

type IndexConversationMetadata struct {
	Next *string `json:"next"`
	Prev *string `json:"prev"`
}

func NewIndexConversationResponse(hasNext bool, inputPagination pagination.Pagination, conversations []entities.ConversationElement) IndexConversationResponse {
	return IndexConversationResponse{
		Metadata: IndexConversationMetadata{
			Prev: inputPagination.Prev(),
			Next: inputPagination.Next(hasNext),
		},
		Data: conversations,
	}
}
