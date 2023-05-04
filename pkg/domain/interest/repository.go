package interest

import (
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

type Repository interface {
	InsertBio(interestEntities.Bio) (string, error)
	GetBioByUserId(string) (interestEntities.Bio, error)
	UpdateBio(interestEntities.Bio) error

	InsertHobbiesByInterestId(string, []interestEntities.Hobbie) error
	GetHobbiesByInterestId(string) ([]interestEntities.Hobbie, error)
	UpdateHobbiesByInterestId(string, []interestEntities.Hobbie) error
	DeleteHobbiesByInterestId(string, []string) error

	InsertMovieSeriesByInterestId(string, []interestEntities.MovieSerie) error
	GetMovieSeriesByInterestId(string) ([]interestEntities.MovieSerie, error)
	UpdateMovieSeriesByInterestId(string, []interestEntities.MovieSerie) error
	DeleteMovieSeriesByInterestId(string, []string) error

	InsertTravelingByInterestId(string, []interestEntities.Travel) error
	GetTravelingByInterestId(string) ([]interestEntities.Travel, error)
	UpdateTravelingByInterestId(string, []interestEntities.Travel) error
	DeleteTravelingByInterestId(string, []string) error

	InsertSportByInterestId(string, []interestEntities.Sport) error
	GetSportByInterestId(string) ([]interestEntities.Sport, error)
	UpdateSportByInterestId(string, []interestEntities.Sport) error
	DeleteSportByInterestId(string, []string) error
}
