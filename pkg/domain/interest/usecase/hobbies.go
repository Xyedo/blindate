package usecase

import (
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

// CreateHobbiesByInterestId implements interest.Usecase
func (i *interestUC) CreateHobbiesByInterestId(interestId string, hobbies []string) ([]string, error) {
	hobbiesDb := make([]interestEntities.Hobbie, 0, len(hobbies))
	for _, hobbie := range hobbies {
		hobbiesDb = append(hobbiesDb, interestEntities.Hobbie{
			Hobbie: hobbie,
		})
	}

	err := i.interestRepo.InsertHobbiesByInterestId(interestId, hobbiesDb)
	if err != nil {
		return nil, err
	}

	returnedIds := make([]string, 0, len(hobbiesDb))
	for _, hobbieDb := range hobbiesDb {
		returnedIds = append(returnedIds, hobbieDb.Id)
	}

	return returnedIds, nil
}

// UpdateHobbiesByInterestId implements interest.Usecase
func (i *interestUC) UpdateHobbiesByInterestId(id string, hobbies []interestDTOs.Hobbie) error {
	hobbieEntities := make([]interestEntities.Hobbie, 0, len(hobbies))
	for _, hobbie := range hobbies {
		hobbieEntities = append(hobbieEntities, interestEntities.Hobbie(hobbie))
	}

	err := i.interestRepo.UpdateHobbiesByInterestId(id, hobbieEntities)
	if err != nil {
		return err
	}

	return nil
}

// DeleteHobbiesByInterestId implements interest.Usecase
func (i *interestUC) DeleteHobbiesByInterestId(id string, hobbieIds []string) error {
	return i.interestRepo.DeleteHobbiesByInterestId(id, hobbieIds)
}
