package dtos

import "github.com/invopop/validation"

type IndexMatchsQueryParams struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (params *IndexMatchsQueryParams) Mod() *IndexMatchsQueryParams {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 10
	}

	return params
}

func (params IndexMatchsQueryParams) Validate() error {
	return validation.ValidateStruct(&params,
		validation.Field(&params.Page, validation.Required, validation.Min(1)),
		validation.Field(&params.Limit, validation.Required, validation.Min(1)),
	)
}

type PutTransitionRequest struct {
	Swipe bool `json:"swipe"`
}
