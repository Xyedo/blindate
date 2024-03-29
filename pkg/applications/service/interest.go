package service

import (
	"strings"

	"github.com/xyedo/blindate/pkg/domain/interest"
	interestEntity "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

func NewInterest(intrRepo interest.Repository) *Interest {
	return &Interest{
		interestRepo: intrRepo,
	}
}

type Interest struct {
	interestRepo interest.Repository
}

func (i *Interest) GetInterest(userId string) (interestEntity.FullDTO, error) {
	intr, err := i.interestRepo.GetInterest(userId)
	if err != nil {
		return interestEntity.FullDTO{}, err
	}
	return intr, nil
}

func (i *Interest) CreateNewBio(intr *interestEntity.BioDTO) error {
	intr.Bio = strings.TrimSpace(intr.Bio)
	err := i.interestRepo.InsertInterestBio(intr)
	if err != nil {
		return err
	}
	err = i.interestRepo.InsertNewStats(intr.Id)
	if err != nil {
		return err
	}
	return nil
}
func (i *Interest) GetBio(userId string) (interestEntity.BioDTO, error) {
	bio, err := i.interestRepo.SelectInterestBio(userId)
	if err != nil {

		return interestEntity.BioDTO{}, err
	}
	return bio, nil
}

func (i *Interest) PutBio(bio interestEntity.BioDTO) error {
	err := i.interestRepo.UpdateInterestBio(bio)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interest) CreateNewHobbies(interestId string, hobbies []string) ([]interestEntity.HobbieDTO, error) {
	hobbiesDTO := make([]interestEntity.HobbieDTO, 0, len(hobbies))
	for _, hobbie := range hobbies {
		hobbiesDTO = append(hobbiesDTO, interestEntity.HobbieDTO{
			Hobbie: hobbie,
		})
	}
	err := i.interestRepo.InsertInterestHobbies(interestId, hobbiesDTO)
	if err != nil {
		return nil, err
	}
	return hobbiesDTO, nil
}

func (i *Interest) PutHobbies(interestId string, hobbies []interestEntity.HobbieDTO) error {
	err := i.interestRepo.UpdateInterestHobbies(interestId, hobbies)
	if err != nil {
		return err
	}

	return nil
}

func (i *Interest) DeleteHobbies(interestId string, ids []string) ([]string, error) {
	deletedIds, err := i.interestRepo.DeleteInterestHobbies(interestId, ids)
	if err != nil {
		return nil, err
	}

	return deletedIds, nil
}

func (i *Interest) CreateNewMovieSeries(interestId string, movieSeries []string) ([]interestEntity.MovieSerieDTO, error) {
	movieSeriesDTO := make([]interestEntity.MovieSerieDTO, 0, len(movieSeries))
	for _, movieSerie := range movieSeries {
		movieSeriesDTO = append(movieSeriesDTO, interestEntity.MovieSerieDTO{
			MovieSerie: movieSerie,
		})
	}
	err := i.interestRepo.InsertInterestMovieSeries(interestId, movieSeriesDTO)

	if err != nil {
		return nil, err
	}
	return movieSeriesDTO, nil
}

func (i *Interest) PutMovieSeries(interestId string, movieSeries []interestEntity.MovieSerieDTO) error {
	err := i.interestRepo.UpdateInterestMovieSeries(interestId, movieSeries)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interest) DeleteMovieSeries(interestId string, ids []string) ([]string, error) {
	deletedIds, err := i.interestRepo.DeleteInterestMovieSeries(interestId, ids)
	if err != nil {
		return nil, err
	}
	return deletedIds, nil
}

func (i *Interest) CreateNewTraveling(interestId string, travels []string) ([]interestEntity.TravelDTO, error) {
	travelsDTO := make([]interestEntity.TravelDTO, 0, len(travels))
	for _, travel := range travels {
		travelsDTO = append(travelsDTO, interestEntity.TravelDTO{
			Travel: travel,
		})
	}
	err := i.interestRepo.InsertInterestTraveling(interestId, travelsDTO)
	if err != nil {
		return nil, err
	}
	return travelsDTO, nil
}
func (i *Interest) PutTraveling(interestId string, travels []interestEntity.TravelDTO) error {
	err := i.interestRepo.UpdateInterestTraveling(interestId, travels)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interest) DeleteTravels(interestId string, ids []string) ([]string, error) {
	deletedIds, err := i.interestRepo.DeleteInterestTraveling(interestId, ids)
	if err != nil {
		return nil, err
	}
	return deletedIds, nil
}
func (i *Interest) CreateNewSports(interestId string, sports []string) ([]interestEntity.SportDTO, error) {
	sportDTO := make([]interestEntity.SportDTO, 0, len(sports))
	for _, sport := range sports {
		sportDTO = append(sportDTO, interestEntity.SportDTO{
			Sport: sport,
		})
	}
	err := i.interestRepo.InsertInterestSports(interestId, sportDTO)
	if err != nil {
		return nil, err
	}
	return sportDTO, nil
}
func (i *Interest) PutSports(interestId string, sports []interestEntity.SportDTO) error {
	err := i.interestRepo.UpdateInterestSport(interestId, sports)
	if err != nil {
		return err
	}

	return nil
}

func (i *Interest) DeleteSports(interestId string, ids []string) ([]string, error) {
	deletedIds, err := i.interestRepo.DeleteInterestSports(interestId, ids)
	if err != nil {
		return nil, err
	}
	return deletedIds, nil
}
