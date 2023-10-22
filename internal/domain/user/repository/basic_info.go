package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func StoreBasicInfo(ctx context.Context, conn pg.Querier, payload entities.BasicInfo) (string, error) {
	const storeBasicInfo = `
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
			updated_at,
			version
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15
		)
		returning user_id
	`
	var returningUserId string
	err := conn.
		QueryRow(ctx, storeBasicInfo,
			payload.UserId,
			payload.Gender,
			payload.FromLoc,
			payload.Height,
			payload.EducationLevel,
			payload.Drinking,
			payload.Smoking,
			payload.RelationshipPref,
			payload.LookingFor,
			payload.Zodiac,
			payload.Kids,
			payload.Work,
			payload.CreatedAt,
			payload.UpdatedAt,
			payload.Version,
		).Scan(&returningUserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			panic("invalid store")
		}

		return "", err
	}

	return returningUserId, nil

}

func GetBasicInfoById(ctx context.Context, conn pg.Querier, id string, opts ...entities.GetBasicInfoOption) (entities.BasicInfo, error) {
	const getBasicInfoById = `
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
		updated_at,
		version
	FROM basic_info
	WHERE user_id = $1
`

	query := getBasicInfoById
	if len(opts) > 0 && opts[0].PessimisticLocking {
		query += "\nSELECT FOR UPDATE"
	}

	var returnedBasicInfo entities.BasicInfo
	err := conn.QueryRow(ctx, query, id).
		Scan(
			&returnedBasicInfo.UserId,
			&returnedBasicInfo.Gender,
			&returnedBasicInfo.FromLoc,
			&returnedBasicInfo.Height,
			&returnedBasicInfo.EducationLevel,
			&returnedBasicInfo.Drinking,
			&returnedBasicInfo.Smoking,
			&returnedBasicInfo.RelationshipPref,
			&returnedBasicInfo.LookingFor,
			&returnedBasicInfo.Zodiac,
			&returnedBasicInfo.Kids,
			&returnedBasicInfo.Work,
			&returnedBasicInfo.CreatedAt,
			&returnedBasicInfo.UpdatedAt,
			&returnedBasicInfo.Version,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.BasicInfo{}, apperror.NotFound(apperror.Payload{
				Error: err,
			})
		}

		return entities.BasicInfo{}, err
	}

	return returnedBasicInfo, nil

}

func UpdateBasicInfoById(ctx context.Context, conn pg.Querier, id string, payload entities.UpdateBasicInfo) error {
	const updateBasicInfoById = `
	UPDATE basic_info SET 
		gender = CASE WHEN $1 THEN $2 ELSE gender,
		from_loc = CASE WHEN $3 THEN $4 ELSE from_loc END,
		height = CASE WHEN $5 THEN $6 ELSE height END,
		education_level =  CASE WHEN $7 THEN $8 ELSE education_level END,
		drinking =  CASE WHEN $9 THEN $10 ELSE drinking END,
		smoking =  CASE WHEN $11 THEN $12 ELSE smoking END,
		relationship_pref =  CASE WHEN $13 THEN $14 ELSE relationship_pref END,
		looking_for =  CASE WHEN $15 THEN $16 ELSE looking_for END,
		zodiac =  CASE WHEN $17 THEN $18 ELSE zodiac END,
		kids =  CASE WHEN $19 THEN $20 ELSE kids END,
		work =  CASE WHEN $21 THEN $22 ELSE work END,
		updated_at =  $23,
		version =  version +1
	FROM basic_info
	WHERE user_id = $24
`
	res, err := conn.Exec(ctx, updateBasicInfoById,
		payload.Gender.IsSet(), payload.Gender,
		payload.FromLoc.IsSet(), payload.FromLoc,
		payload.Height.IsSet(), payload.Height,
		payload.EducationLevel.IsSet(), payload.EducationLevel,
		payload.Drinking.IsSet(), payload.Drinking,
		payload.Smoking.IsSet(), payload.Smoking,
		payload.RelationshipPref.IsSet(), payload.RelationshipPref,
		payload.LookingFor.IsSet(), payload.LookingFor,
		payload.Zodiac.IsSet(), payload.Zodiac,
		payload.Kids.IsSet(), payload.Kids,
		payload.Work.IsSet(), payload.Work,
		payload.UpdateAt,
		id,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() != 1 {
		panic("something went wrong")
	}

	return nil
}
