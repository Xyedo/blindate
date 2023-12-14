package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func StoreHobbiesByUserId(ctx context.Context, conn pg.Querier, userId string, hobbies []entities.Hobbie) error {
	if len(hobbies) == 0 {
		return nil
	}
	copyCount, err := conn.CopyFrom(ctx,
		pgx.Identifier{"hobbies"},
		[]string{"account_id", "hobbie", "created_at", "updated_at", "version"},
		pgx.CopyFromSlice(
			len(hobbies),
			func(i int) ([]any, error) {
				return []any{
					userId,
					hobbies[i].Hobbie,
					hobbies[i].CreatedAt,
					hobbies[i].UpdatedAt,
					hobbies[i].Version,
				}, nil
			},
		),
	)
	if err != nil {
		return err
	}

	if copyCount != int64(len(hobbies)) {
		return errors.New("something went wrong")
	}

	return nil
}

func UpdateHobbies(ctx context.Context, conn pg.Querier, hobbies []entities.UpdateHobbie) error {
	if len(hobbies) == 0 {
		return nil
	}

	const updateHobbieById = `
		UPDATE hobbies SET
			hobbie = $1,
			updated_at = $2,
			version = version +1
		WHERE id = $3
	`
	var batch pgx.Batch
	now := time.Now()
	for i := range hobbies {
		batch.Queue(
			updateHobbieById,
			hobbies[i].Hobbie,
			now,
			hobbies[i].Id,
		).Exec(func(ct pgconn.CommandTag) error {
			if ct.RowsAffected() == 0 {
				return errors.New("invalid")
			}
			return nil
		})
	}
	err := conn.SendBatch(ctx, &batch).Close()
	if err != nil {
		return err
	}

	return nil
}

func DeleteHobbiesByIds(ctx context.Context, conn pg.Querier, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	const deleteHobiesByIds = `
		DELETE hobbies WHERE IN (?)
	`
	query, args, err := pg.In(deleteHobiesByIds, ids)
	if err != nil {
		return err
	}

	ct, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() != int64(len(ids)) {
		return errors.New("something went wrong")
	}

	return nil
}

func StoreMovieSeriesByUserId(ctx context.Context, conn pg.Querier, userId string, movieSeries []entities.MovieSerie) error {
	if len(movieSeries) == 0 {
		return nil
	}

	copyCount, err := conn.CopyFrom(ctx,
		pgx.Identifier{"movie_series"},
		[]string{"account_id", "movie_serie", "created_at", "updated_at", "version"},
		pgx.CopyFromSlice(
			len(movieSeries),
			func(i int) ([]any, error) {
				return []any{
					userId,
					movieSeries[i].MovieSerie,
					movieSeries[i].CreatedAt,
					movieSeries[i].UpdatedAt,
					movieSeries[i].Version,
				}, nil
			},
		),
	)
	if err != nil {
		return err
	}

	if copyCount != int64(len(movieSeries)) {
		panic("not match")
	}

	return nil
}

func UpdateMovieSeries(ctx context.Context, conn pg.Querier, movieSeries []entities.UpdateMovieSeries) error {
	if len(movieSeries) == 0 {
		return nil
	}

	const updateMovieSerieById = `
		UPDATE movie_series SET
			movie_serie = $1,
			updated_at = $2,
			version = version +1
		WHERE id = $3
	`
	var batch pgx.Batch
	now := time.Now()
	for i := range movieSeries {
		batch.Queue(
			updateMovieSerieById,
			movieSeries[i].MovieSerie,
			now,
			movieSeries[i].Id,
		)
	}
	br := conn.SendBatch(ctx, &batch)
	defer br.Close()

	ct, err := br.Exec()
	if err != nil {
		return err
	}

	if ct.RowsAffected() != int64(len(movieSeries)) {
		return errors.New("something went wrong")
	}

	return nil
}

func DeleteMovieSeriesByIds(ctx context.Context, conn pg.Querier, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	const deleteMovieSeriesByIds = `
		DELETE movie_series WHERE IN (?)
	`
	query, args, err := pg.In(deleteMovieSeriesByIds, ids)
	if err != nil {
		return err
	}

	ct, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() != int64(len(ids)) {
		return errors.New("something went wrong")
	}

	return nil
}

func StoreTravelingsByUserId(ctx context.Context, conn pg.Querier, userId string, travels []entities.Travel) error {
	if len(travels) == 0 {
		return nil
	}

	copyCount, err := conn.CopyFrom(ctx,
		pgx.Identifier{"traveling"},
		[]string{"account_id", "travel", "created_at", "updated_at", "version"},
		pgx.CopyFromSlice(
			len(travels),
			func(i int) ([]any, error) {
				return []any{
					userId,
					travels[i].Travel,
					travels[i].CreatedAt,
					travels[i].UpdatedAt,
					travels[i].Version,
				}, nil
			},
		),
	)
	if err != nil {
		return err
	}

	if copyCount != int64(len(travels)) {
		return errors.New("something went wrong")
	}

	return nil
}

func UpdateTravelings(ctx context.Context, conn pg.Querier, travels []entities.UpdateTravel) error {
	if len(travels) == 0 {
		return nil
	}

	const upsertTravelingById = `
		UPDATE traveling SET
			travel = $1,
			updated_at = $2,
			version = version +1
		WHERE id = $3
	`
	var batch pgx.Batch
	now := time.Now()
	for i := range travels {
		batch.Queue(
			upsertTravelingById,
			travels[i].Travel,
			now,
			travels[i].Id,
		)
	}
	br := conn.SendBatch(ctx, &batch)
	defer br.Close()

	ct, err := br.Exec()
	if err != nil {
		return err
	}

	if ct.RowsAffected() != int64(len(travels)) {
		return errors.New("something went wrong")
	}

	return nil
}

func DeleteTravelingByIds(ctx context.Context, conn pg.Querier, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	const deleteTravelingByIds = `
		DELETE traveling WHERE IN (?)
	`
	query, args, err := pg.In(deleteTravelingByIds, ids)
	if err != nil {
		return err
	}

	ct, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() != int64(len(ids)) {
		return errors.New("something went wrong")
	}

	return nil
}

func StoreSportsByUserId(ctx context.Context, conn pg.Querier, userId string, sports []entities.Sport) error {
	if len(sports) == 0 {
		return nil
	}

	copyCount, err := conn.CopyFrom(ctx,
		pgx.Identifier{"sports"},
		[]string{"account_id", "sport", "created_at", "updated_at", "version"},
		pgx.CopyFromSlice(
			len(sports),
			func(i int) ([]any, error) {
				return []any{
					userId,
					sports[i].Sport,
					sports[i].CreatedAt,
					sports[i].UpdatedAt,
					1,
				}, nil
			},
		),
	)
	if err != nil {
		return err
	}

	if copyCount != int64(len(sports)) {
		return errors.New("something went wrong")
	}

	return nil
}

func UpdateSports(ctx context.Context, conn pg.Querier, sports []entities.UpdateSport) error {
	if len(sports) == 0 {
		return nil
	}

	const upsertTravelingById = `
		UPDATE sports SET
			sport = $1,
			updated_at = $2,
			version = version +1
		WHERE id = $3
	`
	var batch pgx.Batch
	now := time.Now()
	for i := range sports {
		batch.Queue(
			upsertTravelingById,
			sports[i].Sport,
			now,
			sports[i].Id,
		)
	}
	br := conn.SendBatch(ctx, &batch)
	defer br.Close()

	ct, err := br.Exec()
	if err != nil {
		return err
	}

	if ct.RowsAffected() != int64(len(sports)) {
		return errors.New("something went wrong")
	}

	return nil
}

func DeleteSportByIds(ctx context.Context, conn pg.Querier, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	const deleteSportByIds = `
		DELETE sports WHERE IN (?)
	`
	query, args, err := pg.In(deleteSportByIds, ids)
	if err != nil {
		return err
	}

	ct, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() != int64(len(ids)) {
		return errors.New("something went wrong")
	}

	return nil
}
