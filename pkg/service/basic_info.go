package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
	"github.com/xyedo/blindate/pkg/repository"
)

var (
	ErrRefUserIdField           = fmt.Errorf("%w::user_id", domain.ErrRefNotFound23503)
	ErrRefGenderField           = fmt.Errorf("%w::gender", domain.ErrRefNotFound23503)
	ErrRefEducationLevelField   = fmt.Errorf("%w::education", domain.ErrRefNotFound23503)
	ErrRefDrinkingField         = fmt.Errorf("%w::drinking", domain.ErrRefNotFound23503)
	ErrRefSmokingField          = fmt.Errorf("%w::smoking", domain.ErrRefNotFound23503)
	ErrRefRelationshipPrefField = fmt.Errorf("%w::relationship_pref", domain.ErrRefNotFound23503)
	ErrRefLookingForField       = fmt.Errorf("%w::looking_for", domain.ErrRefNotFound23503)
	ErrRefZodiacField           = fmt.Errorf("%w::zodiac", domain.ErrRefNotFound23503)
)

func NewBasicInfo(bInfoRepo repository.BasicInfoRepo) *basicInfo {
	return &basicInfo{
		BasicInfoRepo: bInfoRepo,
	}
}

type basicInfo struct {
	BasicInfoRepo repository.BasicInfoRepo
}

func (b *basicInfo) CreateBasicInfo(bInfo *domain.BasicInfo) error {
	rows, err := b.BasicInfoRepo.InsertBasicInfo(b.domainToEntity(bInfo))
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		if err := parsingPostgreError(err); err != nil {
			return err
		}
		return err
	}
	if rows == 0 {
		panic("rows affected should not be zero")
	}
	return nil
}

func (b *basicInfo) GetBasicInfoByUserId(id string) (*domain.BasicInfo, error) {
	basicInfo, err := b.BasicInfoRepo.GetBasicInfoByUserId(id)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		return nil, err

	}
	return b.entityToDomain(basicInfo), nil
}

func (b *basicInfo) UpdateBasicInfo(bInfo *domain.BasicInfo) error {
	rows, err := b.BasicInfoRepo.UpdateBasicInfo(b.domainToEntity(bInfo))
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrResourceNotFound
		}
		if err := parsingPostgreError(err); err != nil {
			return err
		}
		return err
	}
	if rows == 0 {
		panic(err)
	}
	return nil
}

func (*basicInfo) entityToDomain(basicInfo *entity.BasicInfo) *domain.BasicInfo {
	return &domain.BasicInfo{
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
func (*basicInfo) domainToEntity(basicInfo *domain.BasicInfo) *entity.BasicInfo {
	return &entity.BasicInfo{
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

func parsingPostgreError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23503" {
			switch {
			case strings.Contains(pqErr.Constraint, "user_id"):
				return ErrRefUserIdField
			case strings.Contains(pqErr.Constraint, "gender"):
				return ErrRefGenderField
			case strings.Contains(pqErr.Constraint, "education_level"):
				return ErrRefEducationLevelField
			case strings.Contains(pqErr.Constraint, "drinking"):
				return ErrRefDrinkingField
			case strings.Contains(pqErr.Constraint, "smoking"):
				return ErrRefSmokingField
			case strings.Contains(pqErr.Constraint, "relationship_pref"):
				return ErrRefRelationshipPrefField
			case strings.Contains(pqErr.Constraint, "looking_for"):
				return ErrRefLookingForField
			case strings.Contains(pqErr.Constraint, "zodiac"):
				return ErrRefZodiacField
			}
		}
		if pqErr.Code == "23505" {
			return ErrUniqueConstrainUserId
		}
		return pqErr
	}
	return nil
}
