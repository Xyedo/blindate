package entities

import (
	"strconv"
	"time"

	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/common/ids"
)

func (i UserDetail) ToHobbieIds() []string {
	ids := make([]string, 0, len(i.Hobbies))
	for _, hobbie := range i.Hobbies {
		ids = append(ids, hobbie.Id)
	}

	return ids
}

func (i UserDetail) ToMovieSerieIds() []string {
	ids := make([]string, 0, len(i.MovieSeries))
	for _, movieSerie := range i.MovieSeries {
		ids = append(ids, movieSerie.Id)
	}

	return ids
}

func (i UserDetail) ToTravelIds() []string {
	ids := make([]string, 0, len(i.Travels))
	for _, travel := range i.Travels {
		ids = append(ids, travel.Id)
	}

	return ids
}

func (i UserDetail) ToSportIds() []string {
	ids := make([]string, 0, len(i.Sports))
	for _, sport := range i.Hobbies {
		ids = append(ids, sport.Id)
	}

	return ids
}

type CreateInterest struct {
	Hobbies     []string
	MovieSeries []string
	Travels     []string
	Sports      []string
}

func (payload CreateInterest) Validate(userDetailDb UserDetail) error {
	errUniquePayloads := make(map[string][]string)
	errTooMuchPayload := make(map[string][]string)
	const valueNotUnique = "value is not unique"
	const valueTooMuch = "value is exceeding maximum value"

	if len(payload.Hobbies)+len(userDetailDb.Hobbies) >= 10 {
		errTooMuchPayload["hobbies"] = append(errTooMuchPayload["hobbies"], valueTooMuch)
	}
	if len(payload.MovieSeries)+len(userDetailDb.MovieSeries) >= 10 {
		errTooMuchPayload["movie_series"] = append(errTooMuchPayload["movie_series"], valueTooMuch)
	}
	if len(payload.Travels)+len(userDetailDb.Travels) >= 10 {
		errTooMuchPayload["travels"] = append(errTooMuchPayload["travels"], valueTooMuch)
	}
	if len(payload.Sports)+len(userDetailDb.Sports) >= 10 {
		errTooMuchPayload["sports"] = append(errTooMuchPayload["sports"], valueTooMuch)
	}

	uniqueHobbies := make(map[string]struct{})
	for i := range userDetailDb.Hobbies {
		uniqueHobbies[userDetailDb.Hobbies[i].Hobbie] = struct{}{}
	}

	uniqueMovieSeries := make(map[string]struct{})
	for i := range userDetailDb.MovieSeries {
		uniqueMovieSeries[userDetailDb.MovieSeries[i].MovieSerie] = struct{}{}
	}

	uniqueTravelings := make(map[string]struct{})
	for i := range userDetailDb.Travels {
		uniqueTravelings[userDetailDb.Travels[i].Travel] = struct{}{}
	}

	uniqueSports := make(map[string]struct{})
	for i := range userDetailDb.Sports {
		uniqueSports[userDetailDb.Sports[i].Sport] = struct{}{}
	}

	for i, payloadHobbie := range payload.Hobbies {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueHobbies[payloadHobbie]; ok {
			errUniquePayloads["hobbies."+iStr] = append(errUniquePayloads["hobbies."+iStr], valueNotUnique)
		}
	}

	for i, payloadMovieSerie := range payload.MovieSeries {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueMovieSeries[payloadMovieSerie]; ok {
			errUniquePayloads["movie_series."+iStr] = append(errUniquePayloads["movie_series."+iStr], valueNotUnique)
		}
	}

	for i, payloadTravel := range payload.Travels {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueMovieSeries[payloadTravel]; ok {
			errUniquePayloads["travels."+iStr] = append(errUniquePayloads["travels."+iStr], valueNotUnique)
		}
	}

	for i, payloadSport := range payload.Sports {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueSports[payloadSport]; ok {
			errUniquePayloads["sports."+iStr] = append(errUniquePayloads["sports."+iStr], valueNotUnique)
		}
	}

	errPayload := make([]apperror.ErrorPayload, 0, 2)
	if len(errUniquePayloads) > 0 {
		errPayload = append(errPayload, apperror.ErrorPayload{
			Status:  InterestDuplicate,
			Details: errUniquePayloads,
		})
	}

	if len(errTooMuchPayload) > 0 {
		errPayload = append(errPayload, apperror.ErrorPayload{
			Status:  InterestTooLarge,
			Details: errTooMuchPayload,
		})
	}

	if len(errPayload) > 0 {
		return apperror.UnprocessableEntityWithPayloadMap(apperror.PayloadMap{
			Payloads: errPayload,
		})
	}

	return nil
}

func (c CreateInterest) ToHobbies(userId string) []Hobbie {
	now := time.Now()

	res := make([]Hobbie, 0, len(c.Hobbies))
	for _, hobbie := range c.Hobbies {
		res = append(res, Hobbie{
			Id:        ids.Hobbie(),
			UserId:    userId,
			Hobbie:    hobbie,
			CreatedAt: now,
			UpdatedAt: now,
			Version:   1,
		})
	}

	return res
}

func (c CreateInterest) ToMovieSeries(userId string) []MovieSerie {
	now := time.Now()

	res := make([]MovieSerie, 0, len(c.MovieSeries))
	for _, movieSerie := range c.MovieSeries {
		res = append(res, MovieSerie{
			Id:         ids.MovieSerie(),
			UserId:     userId,
			MovieSerie: movieSerie,
			CreatedAt:  now,
			UpdatedAt:  now,
			Version:    1,
		})
	}

	return res
}
func (c CreateInterest) ToTravels(userId string) []Travel {
	now := time.Now()

	res := make([]Travel, 0, len(c.Travels))
	for _, travel := range c.Travels {
		res = append(res, Travel{
			Id:        ids.Travel(),
			UserId:    userId,
			Travel:    travel,
			CreatedAt: now,
			UpdatedAt: now,
			Version:   1,
		})
	}

	return res
}

func (c CreateInterest) ToSports(userId string) []Sport {
	now := time.Now()

	res := make([]Sport, 0, len(c.Sports))
	for _, sport := range c.Sports {
		res = append(res, Sport{
			Id:        ids.Sport(),
			UserId:    userId,
			Sport:     sport,
			CreatedAt: now,
			UpdatedAt: now,
			Version:   1,
		})
	}

	return res
}

type UpdateInterest struct {
	Hobbies     []UpdateHobbie
	MovieSeries []UpdateMovieSeries
	Travels     []UpdateTravel
	Sports      []UpdateSport
}

func (payload UpdateInterest) Validate(userDetailDb UserDetail) error {
	errUniquePayloads := make(map[string][]string)
	errNotFoundPayload := make(map[string][]string)
	const (
		valueNotFound  = "value is not found"
		valueNotUnique = "value is not unique"
	)

	uniqueHobbieIds := make(map[string]struct{})
	uniqueHobbies := make(map[string]struct{})
	for i := range userDetailDb.Hobbies {
		uniqueHobbieIds[userDetailDb.Hobbies[i].Id] = struct{}{}
		uniqueHobbies[userDetailDb.Hobbies[i].Hobbie] = struct{}{}
	}

	uniqueMovieSerieIds := make(map[string]struct{})
	uniqueMovieSeries := make(map[string]struct{})
	for i := range userDetailDb.MovieSeries {
		uniqueMovieSeries[userDetailDb.MovieSeries[i].MovieSerie] = struct{}{}
		uniqueMovieSerieIds[userDetailDb.MovieSeries[i].Id] = struct{}{}
	}

	uniqueTravelingIds := make(map[string]struct{})
	uniqueTravelings := make(map[string]struct{})
	for i := range userDetailDb.Travels {
		uniqueTravelings[userDetailDb.Travels[i].Travel] = struct{}{}
		uniqueTravelingIds[userDetailDb.Travels[i].Id] = struct{}{}
	}

	uniqueSportIds := make(map[string]struct{})
	uniqueSports := make(map[string]struct{})
	for i := range userDetailDb.Sports {
		uniqueSports[userDetailDb.Sports[i].Sport] = struct{}{}
		uniqueSportIds[userDetailDb.Sports[i].Id] = struct{}{}
	}

	for i := range payload.Hobbies {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueHobbieIds[payload.Hobbies[i].Id]; !ok {
			errNotFoundPayload["hobbies."+iStr] = append(errNotFoundPayload["hobbies."+iStr], valueNotFound)

		}
		if _, ok := uniqueHobbies[payload.Hobbies[i].Hobbie]; ok {
			errUniquePayloads["hobbies."+iStr] = append(errUniquePayloads["hobbies."+iStr], valueNotUnique)
		}
	}

	for i := range payload.MovieSeries {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueMovieSerieIds[payload.MovieSeries[i].Id]; !ok {
			errNotFoundPayload["movie_series."+iStr] = append(errNotFoundPayload["movie_series."+iStr], valueNotFound)
		}
		if _, ok := uniqueMovieSeries[payload.MovieSeries[i].MovieSerie]; ok {
			errUniquePayloads["movie_series."+iStr] = append(errUniquePayloads["movie_series."+iStr], valueNotUnique)
		}
	}

	for i := range payload.Travels {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueTravelingIds[payload.Travels[i].Id]; !ok {
			errNotFoundPayload["travels."+iStr] = append(errNotFoundPayload["travels."+iStr], valueNotFound)
		}
		if _, ok := uniqueMovieSeries[payload.Travels[i].Travel]; ok {
			errUniquePayloads["travels."+iStr] = append(errUniquePayloads["travels."+iStr], valueNotUnique)
		}
	}

	for i := range payload.Sports {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueSportIds[payload.Sports[i].Id]; !ok {
			errNotFoundPayload["sports."+iStr] = append(errNotFoundPayload["sports."+iStr], valueNotFound)
		}
		if _, ok := uniqueSports[payload.Sports[i].Sport]; ok {
			errUniquePayloads["sports."+iStr] = append(errUniquePayloads["sports."+iStr], valueNotUnique)
		}
	}

	errPayload := make([]apperror.ErrorPayload, 0, 2)
	if len(errUniquePayloads) > 0 {
		errPayload = append(errPayload, apperror.ErrorPayload{
			Status:  InterestDuplicate,
			Details: errUniquePayloads,
		})
	}

	if len(errNotFoundPayload) > 0 {
		errPayload = append(errPayload, apperror.ErrorPayload{
			Status:  InterestNotFound,
			Details: errNotFoundPayload,
		})
	}

	if len(errPayload) > 0 {
		return apperror.BadPayloadWithPayloadMap(apperror.PayloadMap{
			Payloads: errPayload,
		})
	}

	return nil
}

type UpdateBio struct {
	Bio       string
	UpdatedAt time.Time
}

type UpdateHobbie struct {
	Id     string
	Hobbie string
}

type UpdateMovieSeries struct {
	Id         string
	MovieSerie string
}

type UpdateTravel struct {
	Id     string
	Travel string
}

type UpdateSport struct {
	Id    string
	Sport string
}
type DeleteInterest struct {
	HobbieIds     []string
	MovieSerieIds []string
	TravelIds     []string
	SportIds      []string
}

func (payload DeleteInterest) ValidateIds(userDetailDb UserDetail) error {
	errPayloads := make(map[string][]string)
	const valueNotFound = "value is not found"

	uniqueHobbies := make(map[string]struct{})
	for i := range userDetailDb.Hobbies {
		uniqueHobbies[userDetailDb.Hobbies[i].Id] = struct{}{}
	}

	uniqueMovieSeries := make(map[string]struct{})
	for i := range userDetailDb.MovieSeries {
		uniqueMovieSeries[userDetailDb.MovieSeries[i].Id] = struct{}{}
	}

	uniqueTravelings := make(map[string]struct{})
	for i := range userDetailDb.Travels {
		uniqueTravelings[userDetailDb.Travels[i].Id] = struct{}{}
	}

	uniqueSports := make(map[string]struct{})
	for i := range userDetailDb.Sports {
		uniqueSports[userDetailDb.Sports[i].Id] = struct{}{}
	}

	for i, hobbieId := range payload.HobbieIds {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueHobbies[hobbieId]; ok {
			errPayloads["hobbies."+iStr] = append(errPayloads["hobbies."+iStr], valueNotFound)
		}
	}

	for i, movieSerieId := range payload.MovieSerieIds {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueMovieSeries[movieSerieId]; ok {
			errPayloads["movie_series."+iStr] = append(errPayloads["movie_series."+iStr], valueNotFound)
		}
	}

	for i, travelId := range payload.TravelIds {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueMovieSeries[travelId]; ok {
			errPayloads["travels."+iStr] = append(errPayloads["travels."+iStr], valueNotFound)
		}
	}

	for i, sportId := range payload.SportIds {
		iStr := strconv.Itoa(i)
		if _, ok := uniqueSports[sportId]; ok {
			errPayloads["sports."+iStr] = append(errPayloads["sports."+iStr], valueNotFound)
		}
	}

	if len(errPayloads) > 0 {
		return apperror.UnprocessableEntity(apperror.Payload{})
	}

	return nil
}
