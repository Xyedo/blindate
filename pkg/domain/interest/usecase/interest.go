package usecase

import (
	"github.com/xyedo/blindate/pkg/domain/interest"
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
)

func New(interestRepo interest.Repository) interest.Usecase {
	return &interestUC{
		interestRepo: interestRepo,
	}
}

type interestUC struct {
	interestRepo interest.Repository
}

// GetById implements interest.Usecase
func (i *interestUC) GetById(userId string) (
	interestDTOs.InterestDetail,
	error,
) {
	bio, err := i.GetBioById(userId)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
	}

	hobbiesDb, err := i.interestRepo.GetHobbiesByInterestId(
		bio.InterestId,
	)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
	}

	movieSeriesDb, err := i.interestRepo.GetMovieSeriesByInterestId(
		bio.InterestId,
	)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
	}

	travelsDb, err := i.interestRepo.GetTravelingByInterestId(
		bio.InterestId,
	)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
	}
	
	sportsDb, err := i.interestRepo.GetSportByInterestId(bio.InterestId)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
	}

	hobbies := make([]interestDTOs.Hobbie, 0, len(hobbiesDb))
	for _, hobbieDb := range hobbiesDb {
		hobbies = append(hobbies, interestDTOs.Hobbie(hobbieDb))
	}
	movieSeries := make([]interestDTOs.MovieSerie, 0, len(movieSeriesDb))
	for _, movieSerieDb := range movieSeriesDb {
		movieSeries = append(
			movieSeries,
			interestDTOs.MovieSerie(movieSerieDb),
		)
	}
	travels := make([]interestDTOs.Travel, 0, len(travelsDb))
	for _, travelDb := range travelsDb {
		travels = append(travels, interestDTOs.Travel(travelDb))
	}
	sports := make([]interestDTOs.Sport, 0, len(sportsDb))
	for _, sportDb := range sportsDb {
		sports = append(sports, interestDTOs.Sport(sportDb))
	}

	return interestDTOs.InterestDetail{
		Bio:         bio,
		Hobbies:     hobbies,
		MovieSeries: movieSeries,
		Travels:     travels,
		Sports:      sports,
	}, nil

}
