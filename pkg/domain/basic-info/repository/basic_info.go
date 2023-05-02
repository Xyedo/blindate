package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	basicinfo "github.com/xyedo/blindate/pkg/domain/basic-info"
	basicInfoEntities "github.com/xyedo/blindate/pkg/domain/basic-info/entities"
)

func New(db *sqlx.DB) basicinfo.Repository {
	return &basicInfoConn{
		conn: db,
	}
}

type basicInfoConn struct {
	conn *sqlx.DB
}

// InsertBasicInfo implements basicinfo.Repository
func (b *basicInfoConn) InsertBasicInfo(basicinfo basicInfoEntities.BasicInfo) error {
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

	var returnUserId string
	err := b.conn.GetContext(ctx, &returnUserId, insertBasicInfo, args...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(apperror.Payload{Error: err})
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				switch {
				case strings.Contains(pqErr.Constraint, "user_id"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"user_id": {"value not found"}}})
				case strings.Contains(pqErr.Constraint, "gender"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"gender": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "education_level"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"education_level": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "drinking"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"drinking": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "smoking"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"smoking": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "relationship_pref"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"relationship_pref": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "looking_for"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"looking_for": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "zodiac"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"zodiac": {"invalid enums"}}})
				}
			}
			if pqErr.Code == "23505" {
				return apperror.Conflicted(apperror.Payload{Error: err, Message: "basic info already created"})
			}
		}
		return err
	}
	return nil
}

// GetBasicInfoByUserId implements basicinfo.Repository
func (b *basicInfoConn) GetBasicInfoByUserId(id string) (basicInfoEntities.BasicInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnedBasicInfo basicInfoEntities.BasicInfo
	err := b.conn.GetContext(ctx, &returnedBasicInfo, getBasicInfoByUserId, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return basicInfoEntities.BasicInfo{}, apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return basicInfoEntities.BasicInfo{}, apperror.NotFound(apperror.Payload{Error: err})
		}
		return basicInfoEntities.BasicInfo{}, err
	}
	return returnedBasicInfo, nil
}

// UpdateBasicInfo implements basicinfo.Repository
func (b *basicInfoConn) UpdateBasicInfo(bInfo basicInfoEntities.BasicInfo) error {
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

	var returnedUserId string
	err := b.conn.GetContext(ctx, &returnedUserId, updateBasicInfo, args...)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.NotFound(apperror.Payload{Error: err})
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				switch {
				case strings.Contains(pqErr.Constraint, "user_id"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"user_id": {"value not found"}}})
				case strings.Contains(pqErr.Constraint, "gender"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"gender": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "education_level"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"education_level": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "drinking"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"drinking": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "smoking"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"smoking": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "relationship_pref"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"relationship_pref": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "looking_for"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"looking_for": {"invalid enums"}}})
				case strings.Contains(pqErr.Constraint, "zodiac"):
					return apperror.UnprocessableEntity(apperror.PayloadMap{Error: err, ErrorMap: map[string][]string{"zodiac": {"invalid enums"}}})
				}
			}
		}
		return err
	}

	return nil
}
