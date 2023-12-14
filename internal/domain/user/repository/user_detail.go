package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func StoreUserDetail(ctx context.Context, conn pg.Querier, payload entities.UserDetail) (string, error) {
	const storeUserDetail = `
		INSERT INTO account_detail(
			account_id,
			geog,
			bio,
			last_online,
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
			$1,ST_GeomFromText($2),$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18
		)
		returning account_id
	`
	var returningUserId string
	err := conn.
		QueryRow(ctx, storeUserDetail,
			payload.UserId,
			payload.Geog,
			payload.Bio,
			time.Now(),
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "account_detail_pkey" {
			return "", apperror.Duplicate(
				apperror.Payload{
					Error: err,
				},
				true,
			)
		}

		return "", err
	}

	return returningUserId, nil

}

func GetUserDetailById(ctx context.Context, conn pg.Querier, id string, opts ...entities.GetUserDetailOption) (entities.UserDetail, error) {
	const getUserDetailById = `
	SELECT 
		account_id,
		alias,
		ST_AsText(geog) as geog,
		bio,
		last_online,
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
	FROM account_detail
	WHERE account_id = $1
`

	query := getUserDetailById
	if len(opts) > 0 && opts[0].PessimisticLocking {
		query += "\nFOR UPDATE"
	}
	batch := new(pgx.Batch)

	var returnedUserDetail entities.UserDetail
	batch.Queue(query, id).QueryRow(func(row pgx.Row) error {
		err := row.Scan(
			&returnedUserDetail.UserId,
			&returnedUserDetail.Alias,
			&returnedUserDetail.Geog,
			&returnedUserDetail.Bio,
			&returnedUserDetail.LastOnline,
			&returnedUserDetail.Gender,
			&returnedUserDetail.FromLoc,
			&returnedUserDetail.Height,
			&returnedUserDetail.EducationLevel,
			&returnedUserDetail.Drinking,
			&returnedUserDetail.Smoking,
			&returnedUserDetail.RelationshipPref,
			&returnedUserDetail.LookingFor,
			&returnedUserDetail.Zodiac,
			&returnedUserDetail.Kids,
			&returnedUserDetail.Work,
			&returnedUserDetail.CreatedAt,
			&returnedUserDetail.UpdatedAt,
			&returnedUserDetail.Version,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return apperror.NotFound(apperror.Payload{
					Error:  err,
					Status: entities.UserNotFound,
				})
			}

			return err
		}

		return nil
	})

	if len(opts) > 0 && opts[0].WithHobbies {
		const getHobbieByUserId = `
		SELECT 
			id,
			account_id, 
			hobbie,
			created_at,
			updated_at, 
			version 
		FROM hobbies 
		WHERE account_id = $1`
		batch.Queue(getHobbieByUserId, id).Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var hobbie entities.Hobbie
				err := rows.Scan(
					&hobbie.Id,
					&hobbie.UserId,
					&hobbie.Hobbie,
					&hobbie.CreatedAt,
					&hobbie.UpdatedAt,
					&hobbie.Version,
				)
				if err != nil {
					return err
				}
				returnedUserDetail.Hobbies = append(returnedUserDetail.Hobbies, hobbie)
			}

			return nil
		})
	}

	if len(opts) > 0 && opts[0].WithMovieSeries {
		const getMovieSerieByUserId = `
		SELECT 
			id, 
			account_id,
			movie_serie,
			created_at,
			updated_at, 
			version 
		FROM movie_series 
		WHERE account_id = $1`
		batch.Queue(getMovieSerieByUserId, id).Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var movieSerie entities.MovieSerie
				err := rows.Scan(
					&movieSerie.Id,
					&movieSerie.UserId,
					&movieSerie.MovieSerie,
					&movieSerie.CreatedAt,
					&movieSerie.UpdatedAt,
					&movieSerie.Version,
				)
				if err != nil {
					return err
				}

				returnedUserDetail.MovieSeries = append(returnedUserDetail.MovieSeries, movieSerie)

			}

			return nil
		})

	}

	if len(opts) > 0 && opts[0].WithTravels {
		const getTravelingByUserId = `
		SELECT 
			id, 
			account_id,
			travel,
			created_at,
			updated_at, 
			version 
		FROM traveling 
		WHERE account_id = $1`

		batch.Queue(getTravelingByUserId, id).Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var travel entities.Travel
				err := rows.Scan(
					&travel.Id,
					&travel.UserId,
					&travel.Travel,
					&travel.CreatedAt,
					&travel.UpdatedAt,
					&travel.Version,
				)
				if err != nil {
					return err
				}

				returnedUserDetail.Travels = append(returnedUserDetail.Travels, travel)
			}
			return nil
		})

	}

	if len(opts) > 0 && opts[0].WithSports {
		const getSportByUserId = `
		SELECT 
			id, 
			account_id,
			sport,
			created_at,
			updated_at, 
			version 
		FROM sports 
		WHERE account_id = $1`

		batch.Queue(getSportByUserId, id).Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var sport entities.Sport
				err := rows.Scan(
					&sport.Id,
					&sport.UserId,
					&sport.Sport,
					&sport.CreatedAt,
					&sport.UpdatedAt,
					&sport.Version,
				)
				if err != nil {
					return err
				}

				returnedUserDetail.Sports = append(returnedUserDetail.Sports, sport)
			}

			return nil
		})

	}
	if len(opts) > 0 && opts[0].WithProfilePictures {
		const getPhotoProfile = `
		SELECT 
			id, 
			account_id,
			selected,
			file_id
		FROM profile_pictures 
		WHERE account_id = $1
		ORDER BY selected ASC`

		batch.Queue(getPhotoProfile, id).Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var profilePic entities.ProfilePicture
				err := rows.Scan(
					&profilePic.Id,
					&profilePic.UserId,
					&profilePic.Selected,
					&profilePic.FileId,
				)
				if err != nil {
					return err
				}

				returnedUserDetail.ProfilePictures = append(returnedUserDetail.ProfilePictures, profilePic)
			}

			return nil
		})

	}

	err := conn.SendBatch(ctx, batch).Close()
	if err != nil {
		return entities.UserDetail{}, err
	}

	return returnedUserDetail, nil

}

func UpdateUserDetailById(ctx context.Context, conn pg.Querier, id string, payload entities.UpdateUserDetail) error {
	const updateBasicInfoById = `
	UPDATE account_detail SET 
		gender = CASE WHEN $1 THEN $2 ELSE gender END,
		geog = CASE WHEN $3 THEN ST_GeomFromText($4) ELSE geog END,
		from_loc = CASE WHEN $5 THEN $6 ELSE from_loc END,
		height = CASE WHEN $7 THEN $8 ELSE height END,
		education_level =  CASE WHEN $9 THEN $10 ELSE education_level END,
		drinking =  CASE WHEN $11 THEN $12 ELSE drinking END,
		smoking =  CASE WHEN $13 THEN $14 ELSE smoking END,
		relationship_pref =  CASE WHEN $15 THEN $16 ELSE relationship_pref END,
		looking_for =  CASE WHEN $17 THEN $18 ELSE looking_for END,
		zodiac =  CASE WHEN $19 THEN $20 ELSE zodiac END,
		kids =  CASE WHEN $21 THEN $22 ELSE kids END,
		work =  CASE WHEN $23 THEN $24 ELSE work END,
		updated_at =  $25,
		version =  version +1
	WHERE account_id = $26
`
	res, err := conn.Exec(ctx, updateBasicInfoById,
		payload.Gender.IsSet(), payload.Gender,
		payload.Geog.IsPresent(), payload.Geog.MustGet(),
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
		time.Now(),
		id,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() != 1 {
		return errors.New("something when wrong")
	}

	return nil
}
