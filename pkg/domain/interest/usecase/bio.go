package usecase

import (
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

// CreateBio implements interest.Usecase
func (i *interestUC) CreateBio(bio interestDTOs.Bio) (string, error) {
	bioId, err := i.interestRepo.InsertBio(interestEntities.Bio(bio))
	if err != nil {
		return "", err
	}

	return bioId, nil
}

// GetBioById implements interest.Usecase
func (i *interestUC) GetBioById(userId string) (interestDTOs.Bio, error) {
	bio, err := i.interestRepo.GetBioByUserId(userId)
	if err != nil {
		return interestDTOs.Bio{}, err
	}

	return interestDTOs.Bio(bio), nil
}

// UpdateBio implements interest.Usecase
func (i *interestUC) UpdateBio(bio interestDTOs.UpdateBio) error {
	if !bio.Bio.ValueSet() {
		return apperror.BadPayload(apperror.Payload{
			Message: "body shoud not be emtpy",
		})
	}

	err := i.interestRepo.UpdateBio(interestEntities.Bio{
		InterestId: bio.Id,
		UserId:     bio.UserId,
		Bio:        bio.Bio.MustGet(),
	})
	if err != nil {
		return err
	}

	return nil
}
