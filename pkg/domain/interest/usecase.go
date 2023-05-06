package interest

import (
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
)

type Usecase interface {
	GetById(string) (interestDTOs.InterestDetail, error)

	CreateBio(interestDTOs.Bio) (string, error)
	GetBioById(string) (interestDTOs.Bio, error)
	UpdateBio(interestDTOs.UpdateBio) error

	CreateHobbiesByInterestId(string, []string) ([]string, error)
	UpdateHobbies([]interestDTOs.Hobbie) error
	DeleteHobbiesByIDs([]string) error

	CreateMovieSeriesByInterestId(string, []string) ([]string, error)
	UpdateMovieSeries([]interestDTOs.MovieSerie) error
	DeleteMovieSeriesByIDs([]string) error

	CreateTravelsByInterestId(string, []string) ([]string, error)
	UpdateTravels([]interestDTOs.Travel) error
	DeleteTravelsByIDs([]string) error

	CreateSportsByInterestId(string, []string) ([]string, error)
	UpdateSports([]interestDTOs.Sport) error
	DeleteSportsByIDs([]string) error
}
