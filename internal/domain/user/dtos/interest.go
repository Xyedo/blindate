package dtos

import (
	"strings"

	"github.com/invopop/validation"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/pkg/validator"
)

type PostInterestRequest struct {
	Hobbies     []string `json:"hobbies"`
	MovieSeries []string `json:"movie_series"`
	Travels     []string `json:"travels"`
	Sports      []string `json:"sports"`
}

func (p *PostInterestRequest) Mod() *PostInterestRequest {

	for i := range p.Hobbies {
		p.Hobbies[i] = strings.TrimSpace(p.Hobbies[i])
	}

	for i := range p.MovieSeries {
		p.MovieSeries[i] = strings.TrimSpace(p.MovieSeries[i])
	}

	for i := range p.Travels {
		p.Travels[i] = strings.TrimSpace(p.Travels[i])
	}

	for i := range p.Sports {
		p.Sports[i] = strings.TrimSpace(p.Sports[i])
	}

	return p

}

func (p PostInterestRequest) Validate() error {
	if len(p.Hobbies) == 0 && len(p.MovieSeries) == 0 && len(p.Travels) == 0 && len(p.Sports) == 0 {
		return apperror.BadPayload(apperror.Payload{
			Message: "body should be not empty",
		})
	}
	return validation.ValidateStruct(&p,
		validation.Field(&p.Hobbies, validation.Length(0, 10), validation.By(validator.Unique),
			validation.Each(validation.Required, validation.Length(0, 20)),
		),
		validation.Field(&p.MovieSeries, validation.Length(0, 10), validation.By(validator.Unique),
			validation.Each(validation.Required, validation.Length(0, 20)),
		),
		validation.Field(&p.Travels, validation.Length(0, 10), validation.By(validator.Unique),
			validation.Each(validation.Required, validation.Length(0, 20)),
		),
		validation.Field(&p.Sports, validation.Length(0, 10), validation.By(validator.Unique),
			validation.Each(validation.Required, validation.Length(0, 20)),
		),
	)
}

type UpdateInterest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (u UpdateInterest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Id, validation.Required),
		validation.Field(&u.Name, validation.Required),
	)
}

type PatchInterestRequest struct {
	Hobbies     []UpdateInterest `json:"hobbies"`
	MovieSeries []UpdateInterest `json:"movie_series"`
	Travels     []UpdateInterest `json:"travels"`
	Sports      []UpdateInterest `json:"sports"`
}

func (p *PatchInterestRequest) Mod() *PatchInterestRequest {

	for i := range p.Hobbies {
		p.Hobbies[i].Name = strings.TrimSpace(p.Hobbies[i].Name)
	}

	for i := range p.MovieSeries {
		p.MovieSeries[i].Name = strings.TrimSpace(p.MovieSeries[i].Name)
	}

	for i := range p.Travels {
		p.Travels[i].Name = strings.TrimSpace(p.Travels[i].Name)
	}

	for i := range p.Sports {
		p.Sports[i].Name = strings.TrimSpace(p.Sports[i].Name)
	}

	return p
}

func (p PatchInterestRequest) Validate() error {
	if len(p.Hobbies) == 0 && len(p.MovieSeries) == 0 && len(p.Travels) == 0 && len(p.Sports) == 0 {
		return apperror.BadPayload(apperror.Payload{
			Message: "body should be not empty",
		})
	}
	return validation.ValidateStruct(&p,
		validation.Field(&p.Hobbies,
			validation.Length(1, 10),
			validation.By(validator.UniqueByStructFields(func(i int) any { return p.Hobbies[i].Name })),
		),
		validation.Field(&p.MovieSeries,
			validation.Length(1, 10),
			validation.By(validator.UniqueByStructFields(func(i int) any { return p.MovieSeries[i].Name })),
		),
		validation.Field(&p.Travels,
			validation.Length(1, 10),
			validation.By(validator.UniqueByStructFields(func(i int) any { return p.Travels[i].Name })),
		),
		validation.Field(&p.Sports,
			validation.Length(1, 10),
			validation.By(validator.UniqueByStructFields(func(i int) any { return p.Sports[i].Name })),
		),
	)
}

func (p PatchInterestRequest) ToEntity() entities.UpdateInterest {
	hobbies := make([]entities.UpdateHobbie, 0, len(p.Hobbies))
	for i := range p.Hobbies {
		hobbies = append(hobbies, entities.UpdateHobbie{
			Id:     p.Hobbies[i].Id,
			Hobbie: p.Hobbies[i].Name,
		})
	}

	movieSeries := make([]entities.UpdateMovieSeries, 0, len(p.MovieSeries))
	for i := range p.MovieSeries {
		movieSeries = append(movieSeries, entities.UpdateMovieSeries{
			Id:         p.MovieSeries[i].Id,
			MovieSerie: p.MovieSeries[i].Name,
		})
	}

	travels := make([]entities.UpdateTravel, 0, len(p.Travels))
	for i := range p.Travels {
		travels = append(travels, entities.UpdateTravel{
			Id:     p.Travels[i].Id,
			Travel: p.Travels[i].Name,
		})
	}

	sports := make([]entities.UpdateSport, 0, len(p.Sports))
	for i := range p.Sports {
		sports = append(sports, entities.UpdateSport{
			Id:    p.Sports[i].Id,
			Sport: p.Sports[i].Name,
		})
	}

	return entities.UpdateInterest{
		Hobbies:     hobbies,
		MovieSeries: movieSeries,
		Travels:     travels,
		Sports:      sports,
	}
}

type PostDeleteInterestRequest struct {
	HobbieIds     []string `json:"hobbie_ids"`
	MovieSerieIds []string `json:"movie_serie_ids"`
	TravelIds     []string `json:"travel_ids"`
	SportIds      []string `json:"sport_ids"`
}

func (p PostDeleteInterestRequest) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.HobbieIds, validation.Each(validation.Required)),
		validation.Field(&p.MovieSerieIds, validation.Each(validation.Required)),
		validation.Field(&p.TravelIds, validation.Each(validation.Required)),
		validation.Field(&p.SportIds, validation.Each(validation.Required)),
	)
}
