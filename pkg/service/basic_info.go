package service

import (
	"database/sql"

	"github.com/xyedo/blindate/pkg/domain/basicinfo"
	basicInfoEntity "github.com/xyedo/blindate/pkg/domain/basicinfo/entities"
)

func NewBasicInfo(bInfoRepo basicinfo.Repository) *BasicInfo {
	return &BasicInfo{
		basicInfoRepo: bInfoRepo,
	}
}

type BasicInfo struct {
	basicInfoRepo basicinfo.Repository
}

func (b *BasicInfo) CreateBasicInfo(bInfo basicInfoEntity.FullDTO) error {
	err := b.basicInfoRepo.InsertBasicInfo(b.domainToEntity(bInfo))
	if err != nil {
		return err
	}
	return nil
}

func (b *BasicInfo) GetBasicInfoByUserId(id string) (basicInfoEntity.FullDTO, error) {
	basicInfo, err := b.basicInfoRepo.GetBasicInfoByUserId(id)
	if err != nil {
		return basicInfoEntity.FullDTO{}, err
	}

	return b.entityToDomain(basicInfo), nil
}

func (b *BasicInfo) UpdateBasicInfo(userId string, newBasicInfo basicInfoEntity.Update) error {
	basicInfoDomain, err := b.GetBasicInfoByUserId(userId)
	if err != nil {
		return err
	}
	if newBasicInfo.Gender != nil {
		basicInfoDomain.Gender = *newBasicInfo.Gender
	}
	if newBasicInfo.FromLoc != nil {
		basicInfoDomain.FromLoc = newBasicInfo.FromLoc
	}
	if newBasicInfo.Height != nil {
		basicInfoDomain.Height = newBasicInfo.Height
	}
	if newBasicInfo.EducationLevel != nil {
		basicInfoDomain.EducationLevel = newBasicInfo.EducationLevel
	}
	if newBasicInfo.Drinking != nil {
		basicInfoDomain.Drinking = newBasicInfo.Drinking
	}
	if newBasicInfo.Smoking != nil {
		basicInfoDomain.Smoking = newBasicInfo.Smoking
	}
	if newBasicInfo.RelationshipPref != nil {
		basicInfoDomain.RelationshipPref = newBasicInfo.RelationshipPref
	}
	if newBasicInfo.LookingFor != nil {
		basicInfoDomain.LookingFor = *newBasicInfo.LookingFor
	}
	if newBasicInfo.Zodiac != nil {
		basicInfoDomain.Zodiac = newBasicInfo.Zodiac
	}
	if newBasicInfo.Kids != nil {
		basicInfoDomain.Kids = newBasicInfo.Kids
	}
	if newBasicInfo.Work != nil {
		basicInfoDomain.Work = newBasicInfo.Work
	}

	err = b.basicInfoRepo.UpdateBasicInfo(b.domainToEntity(basicInfoDomain))
	if err != nil {
		return err
	}
	return nil
}

func (BasicInfo) entityToDomain(basicInfo basicInfoEntity.Dao) basicInfoEntity.FullDTO {
	return basicInfoEntity.FullDTO{
		UserId:           basicInfo.UserId,
		Gender:           basicInfo.Gender,
		FromLoc:          newString(basicInfo.FromLoc),
		Height:           newInt(basicInfo.Height),
		EducationLevel:   newString(basicInfo.EducationLevel),
		Drinking:         newString(basicInfo.Drinking),
		Smoking:          newString(basicInfo.Smoking),
		RelationshipPref: newString(basicInfo.RelationshipPref),
		LookingFor:       basicInfo.LookingFor,
		Zodiac:           newString(basicInfo.Zodiac),
		Kids:             newInt(basicInfo.Kids),
		Work:             newString(basicInfo.Work),
		CreatedAt:        basicInfo.CreatedAt,
		UpdatedAt:        basicInfo.UpdatedAt,
	}
}
func (BasicInfo) domainToEntity(basicInfo basicInfoEntity.FullDTO) basicInfoEntity.Dao {
	return basicInfoEntity.Dao{
		UserId:           basicInfo.UserId,
		Gender:           basicInfo.Gender,
		FromLoc:          newNullString(basicInfo.FromLoc),
		Height:           newNullSmallInt(basicInfo.Height),
		EducationLevel:   newNullString(basicInfo.EducationLevel),
		Drinking:         newNullString(basicInfo.Drinking),
		Smoking:          newNullString(basicInfo.Smoking),
		RelationshipPref: newNullString(basicInfo.RelationshipPref),
		LookingFor:       basicInfo.LookingFor,
		Zodiac:           newNullString(basicInfo.Zodiac),
		Kids:             newNullSmallInt(basicInfo.Kids),
		Work:             newNullString(basicInfo.Work),
		CreatedAt:        basicInfo.CreatedAt,
		UpdatedAt:        basicInfo.UpdatedAt,
	}

}
func newNullString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *v,
		Valid:  true,
	}
}

func newNullSmallInt(v *int) sql.NullInt16 {
	if v == nil {
		return sql.NullInt16{}
	}
	return sql.NullInt16{
		Int16: int16(*v),
		Valid: true,
	}
}

func newString(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	return &v.String
}
func newInt(v sql.NullInt16) *int {
	if !v.Valid {
		return nil
	}
	val := int(v.Int16)
	return &val
}
