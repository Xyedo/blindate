package service

import (
	"strings"

	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

func NewInterest(intrRepo repository.Interest) *Interest {
	return &Interest{
		interestRepo: intrRepo,
	}
}

type Interest struct {
	interestRepo repository.Interest
}

func (i *Interest) GetInterest(userId string) (domain.Interest, error) {
	intr, err := i.interestRepo.GetInterest(userId)
	if err != nil {
		return domain.Interest{}, err
	}
	return intr, nil
}

func (i *Interest) CreateNewBio(intr *domain.Bio) error {
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
func (i *Interest) GetBio(userId string) (domain.Bio, error) {
	bio, err := i.interestRepo.SelectInterestBio(userId)
	if err != nil {

		return domain.Bio{}, err
	}
	return bio, nil
}

func (i *Interest) PutBio(bio domain.Bio) error {
	err := i.interestRepo.UpdateInterestBio(bio)
	if err != nil {
		return err
	}
	return nil
}

func (i *Interest) CreateNewHobbies(interestId string, hobbies []string) ([]domain.Hobbie, error) {
	hobbiesDTO := make([]domain.Hobbie, 0, len(hobbies))
	for _, hobbie := range hobbies {
		hobbiesDTO = append(hobbiesDTO, domain.Hobbie{
			Hobbie: hobbie,
		})
	}
	err := i.interestRepo.InsertInterestHobbies(interestId, hobbiesDTO)
	if err != nil {
		return nil, err
	}
	return hobbiesDTO, nil
}

func (i *Interest) PutHobbies(interestId string, hobbies []domain.Hobbie) error {
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

func (i *Interest) CreateNewMovieSeries(interestId string, movieSeries []string) ([]domain.MovieSerie, error) {
	movieSeriesDTO := make([]domain.MovieSerie, 0, len(movieSeries))
	for _, movieSerie := range movieSeries {
		movieSeriesDTO = append(movieSeriesDTO, domain.MovieSerie{
			MovieSerie: movieSerie,
		})
	}
	err := i.interestRepo.InsertInterestMovieSeries(interestId, movieSeriesDTO)

	if err != nil {
		return nil, err
	}
	return movieSeriesDTO, nil
}

func (i *Interest) PutMovieSeries(interestId string, movieSeries []domain.MovieSerie) error {
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

func (i *Interest) CreateNewTraveling(interestId string, travels []string) ([]domain.Travel, error) {
	travelsDTO := make([]domain.Travel, 0, len(travels))
	for _, travel := range travels {
		travelsDTO = append(travelsDTO, domain.Travel{
			Travel: travel,
		})
	}
	err := i.interestRepo.InsertInterestTraveling(interestId, travelsDTO)
	if err != nil {
		return nil, err
	}
	return travelsDTO, nil
}
func (i *Interest) PutTraveling(interestId string, travels []domain.Travel) error {
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
func (i *Interest) CreateNewSports(interestId string, sports []string) ([]domain.Sport, error) {
	sportDTO := make([]domain.Sport, 0, len(sports))
	for _, sport := range sports {
		sportDTO = append(sportDTO, domain.Sport{
			Sport: sport,
		})
	}
	err := i.interestRepo.InsertInterestSports(interestId, sportDTO)
	if err != nil {
		return nil, err
	}
	return sportDTO, nil
}
func (i *Interest) PutSports(interestId string, sports []domain.Sport) error {
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
