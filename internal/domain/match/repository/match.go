package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	matchEntities "github.com/xyedo/blindate/internal/domain/match/entities"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	userEntities "github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func CreateCandidateMatchsById(ctx context.Context, conn pg.Querier, userId string, candidateMatchsIds []string) error {
	rowAffected, err := conn.CopyFrom(ctx,
		pgx.Identifier{"match"},
		[]string{
			"request_from",
			"request_to",
			"request_status",
			"created_at",
			"updated_at",
			"version",
		},
		pgx.CopyFromSlice(len(candidateMatchsIds), func(i int) ([]any, error) {
			return []any{
				userId,
				candidateMatchsIds[i],
				matchEntities.MatchStatusUnknown,
				time.Now(),
				time.Now(),
				1,
			}, nil
		}),
	)

	if err != nil {
		return err
	}

	if rowAffected != int64(len(candidateMatchsIds)) {
		return errors.New("something went wrong")
	}

	return nil

}

func FindUserMatchByStatus(ctx context.Context, conn pg.Querier, userId string, status matchEntities.MatchStatus, limit, page int) ([]matchEntities.MatchUser, error) {
	const findUserMatchByStatus = `
	SELECT 
		m.request_from
		m.request_to
	FROM match m
	JOIN account_detail ad ON
	ad.account_id = m.request_from OR
	ad.account_id = m.request_to
	WHERE 
	ad.account_id = $1 AND
	m.status = $2
	LIMIT $3
	OFFSET $4
	`

	offset := limit*page - limit
	rows, err := conn.Query(ctx, findUserMatchByStatus, userId, string(status), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matchUserIds := make([]string, 0)
	for rows.Next() {
		var (
			requestFrom, requestTo string
		)
		err := rows.Scan(&requestFrom, &requestTo)
		if err != nil {
			return nil, err
		}
		if requestFrom == userId {
			matchUserIds = append(matchUserIds, requestTo)
			continue
		}
		if requestTo == userId {
			matchUserIds = append(matchUserIds, requestFrom)
			continue
		}

		matchUserIds = append(matchUserIds, requestFrom, requestTo)
	}
	const getUserDetailByIds = `
	SELECT 
		account_id,
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
	WHERE account_id IN (?)
`
	const getHobbieByUserIds = `
		SELECT 
			id,
			account_id, 
			hobbie,
			created_at,
			updated_at, 
			version 
		FROM hobbies 
		WHERE account_id IN (?)`
	const getMovieSerieByUserIds = `
		SELECT 
			id, 
			account_id,
			movie_serie,
			created_at,
			updated_at, 
			version 
		FROM movie_series 
		WHERE account_id IN (?)`
	const getTravelingByUserIds = `
		SELECT 
			id, 
			account_id,
			travel,
			created_at,
			updated_at, 
			version 
		FROM traveling 
		WHERE account_id IN (?)`
	const getSportByUserIds = `
		SELECT 
			id, 
			account_id,
			sport,
			created_at,
			updated_at, 
			version 
		FROM sports 
		WHERE account_id IN (?)`
	const getPhotoProfiles = `
		SELECT 
			id, 
			account_id,
			selected,
			file_id
		FROM profile_pictures 
		WHERE account_id IN (?)
		ORDER BY selected ASC`

	q, args, err := pg.In(getUserDetailByIds, matchUserIds)
	if err != nil {
		return nil, err
	}
	rows, err = conn.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	matchUsers := make([]matchEntities.MatchUser, 0, len(matchUserIds))
	userIdToRowIdx := make(map[string]int, len(matchUserIds))
	for rows.Next() {
		var matchUser matchEntities.MatchUser
		err = rows.Scan(
			&matchUser.UserId,
			&matchUser.Geog,
			&matchUser.Bio,
			&matchUser.LastOnline,
			&matchUser.Gender,
			&matchUser.FromLoc,
			&matchUser.Height,
			&matchUser.EducationLevel,
			&matchUser.Drinking,
			&matchUser.Smoking,
			&matchUser.RelationshipPref,
			&matchUser.LookingFor,
			&matchUser.Zodiac,
			&matchUser.Kids,
			&matchUser.Work,
			&matchUser.CreatedAt,
			&matchUser.UpdatedAt,
			&matchUser.Version,
		)
		if err != nil {
			return nil, err
		}
		matchUsers = append(matchUsers, matchUser)
		userIdToRowIdx[matchUser.UserId] = len(matchUsers) - 1
	}
	rows.Close()

	var batch pgx.Batch
	q, args, err = pg.In(getHobbieByUserIds, matchUserIds)
	if err != nil {
		return nil, err
	}

	batch.Queue(q, args...).
		Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var hobbie userEntities.Hobbie
				err := rows.Scan(
					&hobbie.UUID,
					&hobbie.UserId,
					&hobbie.Hobbie,
					&hobbie.CreatedAt,
					&hobbie.UpdatedAt,
					&hobbie.Version,
				)
				if err != nil {
					return err
				}
				idx, ok := userIdToRowIdx[hobbie.UserId]
				if !ok {
					continue
				}

				matchUsers[idx].Hobbies = append(matchUsers[idx].Hobbies, hobbie)
			}
			return nil
		})
	q, args, err = pg.In(getMovieSerieByUserIds, matchUserIds)
	if err != nil {
		return nil, err
	}

	batch.Queue(q, args...).
		Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var movieSerie userEntities.MovieSerie
				err := rows.Scan(
					&movieSerie.UUID,
					&movieSerie.UserId,
					&movieSerie.MovieSerie,
					&movieSerie.CreatedAt,
					&movieSerie.UpdatedAt,
					&movieSerie.Version,
				)
				if err != nil {
					return err
				}
				idx, ok := userIdToRowIdx[movieSerie.UserId]
				if !ok {
					continue
				}

				matchUsers[idx].MovieSeries = append(matchUsers[idx].MovieSeries, movieSerie)
			}
			return nil
		})

	q, args, err = pg.In(getTravelingByUserIds, matchUserIds)
	if err != nil {
		return nil, err
	}

	batch.Queue(q, args...).
		Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var travel userEntities.Travel
				err := rows.Scan(
					&travel.UUID,
					&travel.UserId,
					&travel.Travel,
					&travel.CreatedAt,
					&travel.UpdatedAt,
					&travel.Version,
				)
				if err != nil {
					return err
				}
				idx, ok := userIdToRowIdx[travel.UserId]
				if !ok {
					continue
				}
				matchUsers[idx].Travels = append(matchUsers[idx].Travels, travel)
			}
			return nil
		})
	q, args, err = pg.In(getSportByUserIds, matchUserIds)
	if err != nil {
		return nil, err
	}

	batch.Queue(q, args...).
		Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var sport entities.Sport
				err := rows.Scan(
					&sport.UUID,
					&sport.UserId,
					&sport.Sport,
					&sport.CreatedAt,
					&sport.UpdatedAt,
					&sport.Version,
				)
				if err != nil {
					return err
				}
				idx, ok := userIdToRowIdx[sport.UserId]
				if !ok {
					continue
				}
				matchUsers[idx].Sports = append(matchUsers[idx].Sports, sport)
			}

			return nil
		})
	q, args, err = pg.In(getPhotoProfiles, matchUserIds)
	if err != nil {
		return nil, err
	}

	batch.Queue(q, args...).
		Query(func(rows pgx.Rows) error {
			for rows.Next() {
				var profilePic userEntities.ProfilePicture
				err := rows.Scan(
					&profilePic.UUID,
					&profilePic.UserId,
					&profilePic.Selected,
					&profilePic.FileId,
				)
				if err != nil {
					return err
				}
				idx, ok := userIdToRowIdx[profilePic.UserId]
				if !ok {
					continue
				}
				matchUsers[idx].ProfilePictures = append(matchUsers[idx].ProfilePictures, profilePic)
			}
			return nil
		})

	err = conn.SendBatch(ctx, &batch).Close()
	if err != nil {
		return nil, err
	}

	return matchUsers, nil
}
