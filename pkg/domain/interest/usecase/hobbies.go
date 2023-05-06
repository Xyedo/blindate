package usecase

import (
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

// CreateHobbiesByInterestId implements interest.Usecase
func (i *interestUC) CreateHobbiesByInterestId(
	id string,
	hobbies []string,
) ([]string, error) {
	hobbiesDb := make([]interestEntities.Hobbie, 0, len(hobbies))
	for _, hobbie := range hobbies {
		hobbiesDb = append(hobbiesDb, interestEntities.Hobbie{
			Hobbie: hobbie,
		})
	}

	err := i.interestRepo.CheckInsertHobbiesValid(id, len(hobbiesDb))
	if err != nil {
		return nil, err
	}

	err = i.interestRepo.InsertHobbiesByInterestId(id, hobbiesDb)
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
func (i *interestUC) UpdateHobbies(
	hobbies []interestDTOs.Hobbie,
) error {
	hobbiesEntity := make([]interestEntities.Hobbie, 0, len(hobbies))
	for _, hobbie := range hobbies {
		hobbiesEntity = append(
			hobbiesEntity,
			interestEntities.Hobbie(hobbie),
		)
	}

	err := i.interestRepo.UpdateHobbies(hobbiesEntity)
	if err != nil {
		return err
	}

	return nil
}

// DeleteHobbiesByInterestId implements interest.Usecase
func (i *interestUC) DeleteHobbiesByIDs(hobbieIds []string) error {
	return i.interestRepo.DeleteHobbiesByIDs(hobbieIds)
}
