package usecase

import (
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

// CreateSportsByInterestId implements interest.Usecase
func (i *interestUC) CreateSportsByInterestId(
	id string,
	sports []string,
) ([]string, error) {
	sportsDB := make([]interestEntities.Sport, 0, len(sports))
	for _, sport := range sports {
		sportsDB = append(sportsDB, interestEntities.Sport{
			Sport: sport,
		})
	}

	err := i.interestRepo.CheckInsertSportValid(id, len(sportsDB))
	if err != nil {
		return nil, err
	}

	err = i.interestRepo.InsertSportByInterestId(id, sportsDB)
	if err != nil {
		return nil, err
	}

	returnedIds := make([]string, 0, len(sportsDB))
	for _, sportDB := range sportsDB {
		returnedIds = append(returnedIds, sportDB.Id)
	}

	return returnedIds, nil
}

// UpdateSportsByInterestId implements interest.Usecase
func (i *interestUC) UpdateSports(sports []interestDTOs.Sport) error {
	sportsEntity := make([]interestEntities.Sport, 0, len(sports))
	for _, sport := range sports {
		sportsEntity = append(
			sportsEntity,
			interestEntities.Sport(sport),
		)
	}

	err := i.interestRepo.UpdateSport(sportsEntity)
	if err != nil {
		return err
	}

	return nil
}
// DeleteSportsByInterestId implements interest.Usecase
func (i *interestUC) DeleteSportsByIDs(ids []string) error {
	return i.interestRepo.DeleteSportByIDs(ids)
}


