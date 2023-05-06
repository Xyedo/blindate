package v1

import (
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xyedo/blindate/internal/optional"
	"github.com/xyedo/blindate/pkg/common/mod"
	"github.com/xyedo/blindate/pkg/common/validator"
)

type postBioRequest struct {
	Bio string `json:"bio"`
}

func (b *postBioRequest) mod() *postBioRequest {
	if b == nil {
		return nil
	}
	mod.TrimWhiteSpace(&b.Bio)

	return b
}

func (b postBioRequest) validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Bio, validation.Required, validation.Length(1, 300)),
	)
}

type patchBioRequest struct {
	Bio optional.String `json:"bio"`
}

func (b *patchBioRequest) mod() *patchBioRequest {
	if b == nil {
		return nil
	}
	b.Bio.If(func(bio string) {
		mod.TrimWhiteSpace(&bio)
		b.Bio.Set(bio)
	})

	return b
}

func (b patchBioRequest) validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Bio, validation.Required.When(b.Bio.ValueSet()), validation.Length(1, 300)),
	)
}

type postHobbiesRequest struct {
	Hobies []string `json:"hobbies"`
}

func (h *postHobbiesRequest) mod() *postHobbiesRequest {
	for i, hobbie := range h.Hobies {
		mod.TrimWhiteSpace(&hobbie)
		h.Hobies[i] = hobbie
	}

	return h
}

func (h postHobbiesRequest) validate() error {
	err := validation.ValidateStruct(&h,
		validation.Field(&h.Hobies, validation.Required, validation.Length(1, 10),
			validation.Each(validation.Required, validation.Length(1, 21)),
		),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(h.Hobies, "hobbies")
	if err != nil {
		return err
	}

	return nil
}

type hobbie struct {
	Id     string `json:"id"`
	Hobbie string `json:"hobbie"`
}

func (h hobbie) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Id, validation.Required, is.UUIDv4),
		validation.Field(&h.Hobbie, validation.Required, validation.Length(1, 21)),
	)
}

type patchHobbiesRequest struct {
	Hobies []hobbie `json:"hobbies"`
}

func (hobbies *patchHobbiesRequest) mod() *patchHobbiesRequest {
	for i, hobbie := range hobbies.Hobies {
		mod.TrimWhiteSpace(&hobbie.Hobbie)
		hobbies.Hobies[i].Hobbie = hobbie.Hobbie
	}
	return hobbies
}

func (hobbies patchHobbiesRequest) validate() error {
	err := validation.ValidateStruct(&hobbies,
		validation.Field(&hobbies.Hobies, validation.Required, validation.Length(1, 10)),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(hobbies.Hobies, "hobbies")
	if err != nil {
		return err
	}

	return nil
}

type deleteHobbiesRequest struct {
	IDs []string `json:"hobbie_ids"`
}

func (hobbies *deleteHobbiesRequest) mod() *deleteHobbiesRequest {
	for i := range hobbies.IDs {
		mod.Trim(&hobbies.IDs[i])
	}
	return hobbies
}

func (hobbies deleteHobbiesRequest) validate() error {
	err := validation.ValidateStruct(&hobbies,
		validation.Field(&hobbies.IDs, validation.Required, validation.Length(1, 10),
			validation.Each(validation.Required, is.UUIDv4),
		),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(hobbies.IDs, "hobbies")
	if err != nil {
		return err
	}

	return nil

}

type postMovieSeriesRequest struct {
	MovieSeries []string `json:"movie_series"`
}

func (m *postMovieSeriesRequest) mod() *postMovieSeriesRequest {
	for i, movieSerie := range m.MovieSeries {
		mod.TrimWhiteSpace(&movieSerie)
		m.MovieSeries[i] = movieSerie
	}
	return m
}

func (movieSeries postMovieSeriesRequest) validate() error {
	return validation.ValidateStruct(&movieSeries,
		validation.Field(&movieSeries.MovieSeries, validation.Required, validation.Length(1, 10),
			validation.Each(validation.Required, validation.Length(1, 30)),
		),
	)
}

type movieSerie struct {
	Id         string `json:"id"`
	MovieSerie string `json:"movie_serie"`
}

func (movieSerie movieSerie) Validate() error {
	return validation.ValidateStruct(&movieSerie,
		validation.Field(&movieSerie.Id, validation.Required, is.UUIDv4),
		validation.Field(&movieSerie.MovieSerie, validation.Required, validation.Length(1, 30)),
	)
}

type patchMovieSeriesRequest struct {
	MovieSeries []movieSerie `json:"movie_series"`
}

func (movieSeries *patchMovieSeriesRequest) mod() *patchMovieSeriesRequest {
	for i, movieSerie := range movieSeries.MovieSeries {
		mod.TrimWhiteSpace(&movieSerie.MovieSerie)
		movieSeries.MovieSeries[i].MovieSerie = movieSerie.MovieSerie
	}

	return movieSeries
}

func (movieSeries patchMovieSeriesRequest) validate() error {
	err := validation.ValidateStruct(&movieSeries,
		validation.Field(&movieSeries.MovieSeries, validation.Required, validation.Length(1, 10)),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(movieSeries.MovieSeries, "movie_series")
	if err != nil {
		return err
	}

	return nil
}

type deleteMovieSeriesRequest struct {
	IDs []string `json:"movie_serie_ids"`
}

func (movieSeries *deleteMovieSeriesRequest) mod() *deleteMovieSeriesRequest {
	for i := range movieSeries.IDs {
		mod.Trim(&movieSeries.IDs[i])
	}
	return movieSeries
}

func (movieSeries deleteMovieSeriesRequest) validate() error {
	err := validation.ValidateStruct(&movieSeries,
		validation.Field(&movieSeries.IDs, validation.Required, validation.Length(1, 10),
			validation.Each(validation.Required, is.UUIDv4),
		),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(movieSeries.IDs, "movie_series")
	if err != nil {
		return err
	}

	return nil
}

type postTravelsRequest struct {
	Travels []string `json:"travels"`
}

func (t *postTravelsRequest) mod() *postTravelsRequest {
	for i, travel := range t.Travels {
		mod.TrimWhiteSpace(&travel)
		t.Travels[i] = travel
	}
	return t
}

func (travels postTravelsRequest) validate() error {
	return validation.ValidateStruct(&travels,
		validation.Field(&travels.Travels, validation.Required, validation.Length(1, 10),
			validation.Each(validation.Required, validation.Length(1, 21)),
		),
	)
}

type travel struct {
	Id     string `json:"id"`
	Travel string `json:"travel"`
}

func (t travel) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Id, validation.Required, is.UUIDv4),
		validation.Field(&t.Travel, validation.Required, validation.Length(1, 21)),
	)
}

type patchTravelsRequest struct {
	Travels []travel `json:"travels"`
}

func (travels *patchTravelsRequest) mod() *patchTravelsRequest {
	for i, travel := range travels.Travels {
		mod.TrimWhiteSpace(&travel.Travel)
		travels.Travels[i].Travel = travel.Travel
	}

	return travels
}

func (travels patchTravelsRequest) validate() error {
	err := validation.ValidateStruct(&travels,
		validation.Field(&travels.Travels, validation.Required, validation.Length(1, 10)),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(travels.Travels, "travels")
	if err != nil {
		return err
	}

	return nil
}

type deleteTravelsRequest struct {
	IDs []string `json:"travel_ids"`
}

func (travels *deleteTravelsRequest) mod() *deleteTravelsRequest {
	for i := range travels.IDs {
		mod.Trim(&travels.IDs[i])
	}
	return travels
}

func (travels deleteTravelsRequest) validate() error {
	err := validation.ValidateStruct(&travels,
		validation.Field(&travels.IDs, validation.Required, validation.Length(1, 10),
			validation.Each(validation.Required, is.UUIDv4),
		),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(travels.IDs, "travels")
	if err != nil {
		return err
	}

	return nil
}

type postSportsRequest struct {
	Sports []string `json:"sports"`
}

func (sports *postSportsRequest) mod() *postSportsRequest {
	for i, sport := range sports.Sports {
		mod.TrimWhiteSpace(&sport)
		sports.Sports[i] = sport
	}

	return sports
}

func (sports postSportsRequest) validate() error {
	err := validation.ValidateStruct(&sports,
		validation.Field(&sports.Sports, validation.Required, validation.Length(1, 10),
			validation.Each(validation.Required, validation.Length(1, 21)),
		),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(sports.Sports, "sports")
	if err != nil {
		return err
	}

	return nil
}

type sport struct {
	Id    string `json:"id"`
	Sport string `json:"sport"`
}

func (s sport) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Id, validation.Required, is.UUIDv4),
		validation.Field(&s.Sport, validation.Required, validation.Length(1, 21)),
	)
}

type patchSportsRequest struct {
	Sports []sport `json:"sports"`
}

func (sports *patchSportsRequest) mod() *patchSportsRequest {
	for i, sport := range sports.Sports {
		mod.TrimWhiteSpace(&sport.Sport)
		sports.Sports[i].Sport = sport.Sport
	}

	return sports
}

func (sports patchSportsRequest) validate() error {
	err := validation.ValidateStruct(&sports,
		validation.Field(&sports.Sports, validation.Required, validation.Length(1, 10)),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(sports.Sports, "sports")
	if err != nil {
		return err
	}

	return nil
}

type deleteSportsRequest struct {
	IDs []string `json:"sport_ids"`
}

func (sports *deleteSportsRequest) mod() *deleteSportsRequest {
	for i := range sports.IDs {
		mod.Trim(&sports.IDs[i])
	}
	return sports
}

func (sports deleteSportsRequest) validate() error {
	err := validation.ValidateStruct(&sports,
		validation.Field(&sports.IDs, validation.Required, validation.Length(1, 10),
			validation.Each(validation.Required, is.UUIDv4),
		),
	)
	if err != nil {
		return err
	}

	err = validator.Unique(sports.IDs, "sports")
	if err != nil {
		return err
	}

	return nil
}
