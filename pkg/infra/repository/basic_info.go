package repository

import (
	"context"
	"database/sql"
	"errors"

	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/common"
	basicInfoEntity "github.com/xyedo/blindate/pkg/domain/basicinfo/entities"
)

func NewBasicInfo(db *sqlx.DB) *BInfoConn {
	return &BInfoConn{
		conn: db,
	}
}

type BInfoConn struct {
	conn *sqlx.DB
}

func (b *BInfoConn) InsertBasicInfo(basicinfo basicInfoEntity.Dao) error {
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
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $13)
	RETURNING user_id`
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
	var retUserId string
	err := b.conn.GetContext(ctx, &retUserId, query, args...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return common.WrapError(err, common.ErrResourceNotFound)
		}
		if err := b.parsingPostgreError(err); err != nil {
			return err
		}
		return err
	}
	return nil
}

func (b *BInfoConn) GetBasicInfoByUserId(userId string) (basicInfoEntity.Dao, error) {
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
	var basicinfo basicInfoEntity.Dao
	err := b.conn.GetContext(ctx, &basicinfo, query, userId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return basicInfoEntity.Dao{}, common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return basicInfoEntity.Dao{}, common.WrapError(err, common.ErrResourceNotFound)
		}
		return basicInfoEntity.Dao{}, err
	}
	return basicinfo, nil
}

func (b *BInfoConn) UpdateBasicInfo(bInfo basicInfoEntity.Dao) error {
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
	WHERE user_id = $13
	RETURNING user_id`

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

	var retUserId string
	err := b.conn.GetContext(ctx, &retUserId, query, args...)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return common.WrapError(err, common.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return common.WrapError(err, common.ErrResourceNotFound)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return common.WrapError(err, common.ErrResourceNotFound)
		}
		if err := b.parsingPostgreError(err); err != nil {
			return err
		}
		return err
	}
	return nil
}
func (*BInfoConn) parsingPostgreError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23503" {
			switch {
			case strings.Contains(pqErr.Constraint, "user_id"):
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "userId is invalid")
			case strings.Contains(pqErr.Constraint, "gender"):
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "gender value is not valid enums")
			case strings.Contains(pqErr.Constraint, "education_level"):
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "educationLevel value is not valid enums")
			case strings.Contains(pqErr.Constraint, "drinking"):
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "drinking value is not valid enums")
			case strings.Contains(pqErr.Constraint, "smoking"):
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "smoking value is not valid enums")
			case strings.Contains(pqErr.Constraint, "relationship_pref"):
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "relationshipPref value is not valid enums")
			case strings.Contains(pqErr.Constraint, "looking_for"):
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "lookingFor value is not valid enums")
			case strings.Contains(pqErr.Constraint, "zodiac"):
				return common.WrapErrorWithMsg(err, common.ErrRefNotFound23503, "zodiac value is not valid enums")
			}
		}
		if pqErr.Code == "23505" {
			return common.WrapError(err, common.ErrUniqueConstraint23505)
		}
		return pqErr
	}
	return nil
}
