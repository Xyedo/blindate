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
func (i *interestUC) GetById(userId string) (interestDTOs.InterestDetail, error) {
	bio, err := i.GetBioById(userId)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
	}

	hobbiesDb, err := i.interestRepo.GetHobbiesByInterestId(bio.InterestId)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
	}
	hobbies := make([]interestDTOs.Hobbie, 0, len(hobbiesDb))
	for _, hobbieDb := range hobbiesDb {
		hobbies = append(hobbies, interestDTOs.Hobbie(hobbieDb))
	}

	movieSeriesDb, err := i.interestRepo.GetMovieSeriesByInterestId(bio.InterestId)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
	}
	movieSeries := make([]interestDTOs.MovieSerie, 0, len(movieSeriesDb))
	for _, movieSerieDb := range movieSeriesDb {
		movieSeries = append(movieSeries, interestDTOs.MovieSerie(movieSerieDb))
	}

	travelsDb, err := i.interestRepo.GetTravelingByInterestId(bio.InterestId)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
	}
	travels := make([]interestDTOs.Travel, 0, len(travelsDb))
	for _, travelDb := range travelsDb {
		travels = append(travels, interestDTOs.Travel(travelDb))
	}

	sportsDb, err := i.interestRepo.GetSportByInterestId(bio.InterestId)
	if err != nil {
		return interestDTOs.InterestDetail{}, err
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

// CreateMovieSeriesByInterestId implements interest.Usecase
func (*interestUC) CreateMovieSeriesByInterestId(string, []string) ([]string, error) {
	panic("unimplemented")
}

// CreateSportsByInterestId implements interest.Usecase
func (*interestUC) CreateSportsByInterestId(string, []string) ([]string, error) {
	panic("unimplemented")
}

// CreateTravelsByInterestId implements interest.Usecase
func (*interestUC) CreateTravelsByInterestId(string, []string) ([]string, error) {
	panic("unimplemented")
}

// DeleteMovieSeriesByInterestId implements interest.Usecase
func (*interestUC) DeleteMovieSeriesByInterestId(string, []string) error {
	panic("unimplemented")
}

// DeleteSportsByInterestId implements interest.Usecase
func (*interestUC) DeleteSportsByInterestId(string, []string) error {
	panic("unimplemented")
}

// DeleteTravelsByInterestId implements interest.Usecase
func (*interestUC) DeleteTravelsByInterestId(string, []string) error {
	panic("unimplemented")
}

// UpdateMovieSeriesByInterestId implements interest.Usecase
func (*interestUC) UpdateMovieSeriesByInterestId(string, []interestDTOs.MovieSerie) error {
	panic("unimplemented")
}

// UpdateSportsByInterestId implements interest.Usecase
func (*interestUC) UpdateSportsByInterestId(string, []interestDTOs.Sport) error {
	panic("unimplemented")
}

// UpdateTravelsByInterestId implements interest.Usecase
func (*interestUC) UpdateTravelsByInterestId(string, []interestDTOs.Travel) error {
	panic("unimplemented")
}
