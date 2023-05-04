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
	UpdateHobbiesByInterestId(string, []interestDTOs.Hobbie) error
	DeleteHobbiesByInterestId(string, []string) error

	CreateMovieSeriesByInterestId(string, []string) ([]string, error)
	UpdateMovieSeriesByInterestId(string, []interestDTOs.MovieSerie) error
	DeleteMovieSeriesByInterestId(string, []string) error

	CreateTravelsByInterestId(string, []string) ([]string, error)
	UpdateTravelsByInterestId(string, []interestDTOs.Travel) error
	DeleteTravelsByInterestId(string, []string) error

	CreateSportsByInterestId(string, []string) ([]string, error)
	UpdateSportsByInterestId(string, []interestDTOs.Sport) error
	DeleteSportsByInterestId(string, []string) error
}
