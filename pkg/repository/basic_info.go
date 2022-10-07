package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/entity"
)

type BasicInfo interface {
	InsertBasicInfo(basicinfo *entity.BasicInfo) (int64, error)
	GetBasicInfoByUserId(id string) (*entity.BasicInfo, error)
	UpdateBasicInfo(bInfo *entity.BasicInfo) (int64, error)
}

func NewBasicInfo(db *sqlx.DB) *basicInfo {
	return &basicInfo{
		db,
	}
}

type basicInfo struct {
	*sqlx.DB
}

func (b *basicInfo) InsertBasicInfo(basicinfo *entity.BasicInfo) (int64, error) {
	query := `
	INSERT INTO basic_info(
		user_id, 
		gender, 
		from_loc, 
		height, 
		education_level,
		drinking,
		smoking,
		relationship_pref,
		looking_for, 
		zodiac, 
		kids, 
		work, 
		created_at, 
		updated_at)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $13)`
	args := []any{
		basicinfo.UserId,
		basicinfo.Gender,
		basicinfo.FromLoc,
		basicinfo.Height,
		basicinfo.EducationLevel,
		basicinfo.Drinking,
		basicinfo.Smoking,
		basicinfo.RelationshipPref,
		basicinfo.LookingFor,
		basicinfo.Zodiac,
		basicinfo.Kids,
		basicinfo.Work,
		time.Now(),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := b.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		return 0, err
	}
	return affected, nil
}

func (b *basicInfo) GetBasicInfoByUserId(userId string) (*entity.BasicInfo, error) {
	query := `
		SELECT
			user_id, 
			gender, 
			from_loc, 
			height, 
			education_level, 
			drinking, 
			smoking, 
			relationship_pref, 
			looking_for, 
			zodiac, 
			kids, 
			work, 
			created_at, 
			updated_at
		FROM basic_info
		WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var basicinfo entity.BasicInfo
	err := b.GetContext(ctx, &basicinfo, query, userId)
	if err != nil {
		return nil, err
	}
	return &basicinfo, nil
}

func (b *basicInfo) UpdateBasicInfo(bInfo *entity.BasicInfo) (int64, error) {
	query := `
	UPDATE basic_info SET
		gender =$1, 
		from_loc=$2, 
		height=$3, 
		education_level=$4,
		drinking=$5,
		smoking=$6,
		relationship_pref=$7,
		looking_for=$8, 
		zodiac=$9, 
		kids=$10, 
		work=$11, 
		updated_at=$12
	WHERE user_id = $13`

	args := []any{
		bInfo.Gender,
		bInfo.FromLoc,
		bInfo.Height,
		bInfo.EducationLevel,
		bInfo.Drinking,
		bInfo.Smoking,
		bInfo.RelationshipPref,
		bInfo.LookingFor,
		bInfo.Zodiac,
		bInfo.Kids,
		bInfo.Work,
		time.Now(),
		bInfo.UserId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := b.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil

}
