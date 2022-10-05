package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
	"github.com/xyedo/blindate/pkg/repository"
)

var (
	ErrUserIdField            = errors.New("database: user_id is not valid references")
	ErrGenderField            = errors.New("database: gender is not valid references")
	ErrEducationLevelField    = errors.New("database: education_level is not valid references")
	ErrDrinkingField          = errors.New("database: drinking is not valid references")
	ErrSmokingField           = errors.New("database: smoking is not valid references")
	ErrrRelationshipPrefField = errors.New("database: relationship_pref is not valid references")
	ErrLookingForField        = errors.New("database: looking_for is not valid references")
	ErrZodiacField            = errors.New("database: zodiac is not valid references")
)

type BasicInfo interface {
	CreateBasicInfo(bInfo *domain.BasicInfo) error
	GetBasicInfoByUserId(id string) (*domain.BasicInfo, error)
	UpdateBasicInfo(bInfo *domain.BasicInfo) error
}

func NewBasicInfo(bInfoRepo repository.BasicInfo) *basicInfo {
	return &basicInfo{
		BasicInfoRepo: bInfoRepo,
	}
}

type basicInfo struct {
	BasicInfoRepo repository.BasicInfo
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
		Id:               basicInfo.Id,
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
		Id:               basicInfo.Id,
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
				return ErrUserIdField
			case strings.Contains(pqErr.Constraint, "gender"):
				return ErrGenderField
			case strings.Contains(pqErr.Constraint, "education_level"):
				return ErrEducationLevelField
			case strings.Contains(pqErr.Constraint, "drinking"):
				return ErrDrinkingField
			case strings.Contains(pqErr.Constraint, "smoking"):
				return ErrSmokingField
			case strings.Contains(pqErr.Constraint, "relationship_pref"):
				return ErrrRelationshipPrefField
			case strings.Contains(pqErr.Constraint, "looking_for"):
				return ErrLookingForField
			case strings.Contains(pqErr.Constraint, "zodiac"):
				return ErrZodiacField
			}
		}
		return pqErr
	}
	return nil
}