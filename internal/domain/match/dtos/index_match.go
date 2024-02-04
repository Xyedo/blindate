package dtos

import (
	"github.com/invopop/validation"
	"github.com/xyedo/blindate/internal/domain/match/entities"
	userdtos "github.com/xyedo/blindate/internal/domain/user/dtos"
	"github.com/xyedo/blindate/pkg/pagination"
)

type IndexMatchsQueryParams struct {
	Page   int                       `query:"page"`
	Limit  int                       `query:"limit"`
	Status entities.FilterIndexMatch `query:"status"`
}

func (params *IndexMatchsQueryParams) Mod() *IndexMatchsQueryParams {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Status == "" {
		params.Status = entities.FilterIndexMatchCandidate
	}

	return params
}

func (params IndexMatchsQueryParams) Validate() error {
	return validation.ValidateStruct(&params,
		validation.Field(&params.Page, validation.Required, validation.Min(1)),
		validation.Field(&params.Limit, validation.Required, validation.Min(1)),
		validation.Field(&params.Status,
			validation.In(
				entities.FilterIndexMatchCandidate,
				entities.FilterIndexMatchAccepted,
				entities.FilterIndexMatchLikes,
			),
		),
	)
}

type IndexMatchResponse struct {
	Metadata IndexMatchResponseMetadata `json:"metadata"`
	Data     []IndexMatchElement        `json:"data"`
}

type IndexMatchResponseMetadata struct {
	Prev *string `json:"prev"`
	Next *string `json:"next"`
}

type IndexMatchElement struct {
	Id       string  `json:"id"`
	Distance float64 `json:"distance"`
	userdtos.UserDetail
}

func NewIndexMatchResponse(hasNext bool, inputPagination pagination.Pagination, matchUsers []entities.MatchUser) IndexMatchResponse {
	indexMatch := make([]IndexMatchElement, 0, len(matchUsers))
	for _, matchUser := range matchUsers {
		userDetail := userdtos.NewUserDetailResponse(matchUser.UserDetail)
		indexMatch = append(indexMatch, IndexMatchElement{
			Id:         matchUser.MatchId,
			Distance:   matchUser.Distance,
			UserDetail: userDetail,
		})
	}

	return IndexMatchResponse{
		Metadata: IndexMatchResponseMetadata{
			Prev: inputPagination.Prev(),
			Next: inputPagination.Next(hasNext),
		},
		Data: indexMatch,
	}

}
