package interest

import (
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

type Repository interface {
	InsertBio(interestEntities.Bio) (string, error)
	GetBioByUserId(string) (interestEntities.Bio, error)
	UpdateBio(interestEntities.Bio) error

	CheckInsertHobbiesValid(string, int) error
	InsertHobbiesByInterestId(string, []interestEntities.Hobbie) error
	GetHobbiesByInterestId(string) ([]interestEntities.Hobbie, error)
	UpdateHobbies([]interestEntities.Hobbie) error
	DeleteHobbiesByIDs([]string) error

	CheckInsertMovieSeriesValid(string, int) error
	InsertMovieSeriesByInterestId(string, []interestEntities.MovieSerie) error
	GetMovieSeriesByInterestId(string) ([]interestEntities.MovieSerie, error)
	UpdateMovieSeries([]interestEntities.MovieSerie) error
	DeleteMovieSeriesByIDs([]string) error

	CheckInsertTravelingValid(string, int) error
	InsertTravelingByInterestId(string, []interestEntities.Travel) error
	GetTravelingByInterestId(string) ([]interestEntities.Travel, error)
	UpdateTraveling([]interestEntities.Travel) error
	DeleteTravelingByIDs([]string) error

	CheckInsertSportValid(string, int) error
	InsertSportByInterestId(string, []interestEntities.Sport) error
	GetSportByInterestId(string) ([]interestEntities.Sport, error)
	UpdateSport([]interestEntities.Sport) error
	DeleteSportByIDs([]string) error
}
