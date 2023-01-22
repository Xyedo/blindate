package interest

import (
	interestEntity "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

type Repository interface {
	InsertNewStats(interestId string) error

	GetInterest(userId string) (interestEntity.FullDTO, error)

	InsertInterestBio(intr *interestEntity.BioDTO) error
	SelectInterestBio(userId string) (interestEntity.BioDTO, error)
	UpdateInterestBio(intr interestEntity.BioDTO) error

	InsertInterestHobbies(interestId string, hobbies []interestEntity.HobbieDTO) error
	UpdateInterestHobbies(interestId string, hobbies []interestEntity.HobbieDTO) error
	DeleteInterestHobbies(interestId string, ids []string) ([]string, error)

	InsertInterestMovieSeries(interestId string, movieSeries []interestEntity.MovieSerieDTO) error
	UpdateInterestMovieSeries(interestId string, movieSeries []interestEntity.MovieSerieDTO) error
	DeleteInterestMovieSeries(interestId string, ids []string) ([]string, error)

	InsertInterestTraveling(interestId string, travels []interestEntity.TravelDTO) error
	UpdateInterestTraveling(interestId string, travels []interestEntity.TravelDTO) error
	DeleteInterestTraveling(interestId string, ids []string) ([]string, error)

	InsertInterestSports(interestId string, sports []interestEntity.SportDTO) error
	UpdateInterestSport(interestId string, sports []interestEntity.SportDTO) error
	DeleteInterestSports(interestId string, ids []string) ([]string, error)
}
