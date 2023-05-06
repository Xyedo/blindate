package usecase

import (
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

// CreateTravelsByInterestId implements interest.Usecase
func (i *interestUC) CreateTravelsByInterestId(id string, travels []string) ([]string, error) {
	travelsDB := make([]interestEntities.Travel, 0, len(travels))
	for _, travel := range travels {
		travelsDB = append(travelsDB, interestEntities.Travel{
			Travel: travel,
		})
	}

	err := i.interestRepo.CheckInsertTravelingValid(id, len(travelsDB))
	if err != nil {
		return nil, err
	}

	err = i.interestRepo.InsertTravelingByInterestId(id, travelsDB)
	if err != nil {
		return nil, err
	}

	returnedIds := make([]string, 0, len(travelsDB))
	for _, travelDB := range travelsDB {
		returnedIds = append(returnedIds, travelDB.Id)
	}

	return returnedIds, nil
}

// UpdateTravelsByInterestId implements interest.Usecase
func (i *interestUC) UpdateTravels(travels []interestDTOs.Travel) error {
	travelsEntity := make([]interestEntities.Travel, 0, len(travels))
	for _, travel := range travels {
		travelsEntity = append(
			travelsEntity,
			interestEntities.Travel(travel),
		)
	}

	err := i.interestRepo.UpdateTraveling(travelsEntity)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTravelsByInterestId implements interest.Usecase
func (i *interestUC) DeleteTravelsByIDs(ids []string) error {
	return i.interestRepo.DeleteTravelingByIDs(ids)
}
