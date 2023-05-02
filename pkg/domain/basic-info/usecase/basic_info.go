package usecase

import (
	basicinfo "github.com/xyedo/blindate/pkg/domain/basic-info"
	basicInfoDTOs "github.com/xyedo/blindate/pkg/domain/basic-info/dtos"
	basicInfoEntities "github.com/xyedo/blindate/pkg/domain/basic-info/entities"
)

func New(basicInfoRepo basicinfo.Repository) basicinfo.Usecase {
	return &basicInfoUC{
		basicInfoRepo: basicInfoRepo,
	}
}

type basicInfoUC struct {
	basicInfoRepo basicinfo.Repository
}

// Create implements basicinfo.Usecase
func (b *basicInfoUC) Create(basicInfo basicInfoDTOs.CreateBasicInfo) error {
	err := b.basicInfoRepo.InsertBasicInfo(basicInfoEntities.BasicInfo(basicInfo))
	if err != nil {
		return err
	}

	return nil
}

// GetById implements basicinfo.Usecase
func (b *basicInfoUC) GetById(id string) (basicInfoEntities.BasicInfo, error) {
	foundBasicInfo, err := b.basicInfoRepo.GetBasicInfoByUserId(id)
	if err != nil {
		return basicInfoEntities.BasicInfo{}, err
	}

	return foundBasicInfo, nil
}

// Update implements basicinfo.Usecase
func (b *basicInfoUC) Update(basicInfo basicInfoDTOs.UpdateBasicInfo) error {
	basicInfoDB, err := b.basicInfoRepo.GetBasicInfoByUserId(basicInfo.UserId)
	if err != nil {
		return err
	}

	if basicInfo.Gender.ValueSet() {
		basicInfoDB.Gender = basicInfo.Gender.MustGet()
	}

	if basicInfo.FromLoc.ValueSet() {
		basicInfoDB.FromLoc = basicInfo.FromLoc
	}

	if basicInfo.Height.ValueSet() {
		basicInfoDB.Height = basicInfo.Height
	}

	if basicInfo.EducationLevel.ValueSet() {
		basicInfoDB.EducationLevel = basicInfo.EducationLevel
	}

	if basicInfo.Drinking.ValueSet() {
		basicInfoDB.Drinking = basicInfo.Drinking
	}

	if basicInfo.RelationshipPref.ValueSet() {
		basicInfoDB.RelationshipPref = basicInfo.RelationshipPref
	}

	if basicInfo.LookingFor.ValueSet() {
		basicInfoDB.LookingFor = basicInfo.LookingFor.MustGet()
	}

	if basicInfo.Zodiac.ValueSet() {
		basicInfoDB.Zodiac = basicInfo.Zodiac
	}

	if basicInfo.Kids.ValueSet() {
		basicInfoDB.Kids = basicInfo.Kids
	}

	if basicInfo.Work.ValueSet() {
		basicInfoDB.Work = basicInfo.Work
	}

	err = b.basicInfoRepo.UpdateBasicInfo(basicInfoDB)
	if err != nil {
		return err
	}

	return nil
}
