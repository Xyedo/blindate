package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func CreateInterest(ctx context.Context, requestId string, payload entities.CreateInterest) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		userDetail, err := repository.GetUserDetailById(ctx, tx, requestId, entities.GetUserDetailOption{
			PessimisticLocking: true,
			WithHobbies:        true,
			WithMovieSeries:    true,
			WithTravels:        true,
			WithSports:         true,
		})
		if err != nil {
			return err
		}

		err = payload.Validate(userDetail)
		if err != nil {
			return err
		}

		err = repository.StoreHobbiesByUserId(ctx, tx, requestId, payload.ToHobbies(requestId))
		if err != nil {
			return err
		}

		err = repository.StoreMovieSeriesByUserId(ctx, tx, requestId, payload.ToMovieSeries(requestId))
		if err != nil {
			return err
		}

		err = repository.StoreTravelingsByUserId(ctx, tx, requestId, payload.ToTravels(requestId))
		if err != nil {
			return err
		}

		err = repository.StoreSportsByUserId(ctx, tx, requestId, payload.ToSports(requestId))
		if err != nil {
			return err
		}

		return nil
	})

}

func UpdateInterest(ctx context.Context, requestId string, payload entities.UpdateInterest) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		userDetail, err := repository.GetUserDetailById(ctx, tx, requestId, entities.GetUserDetailOption{
			PessimisticLocking: true,
			WithHobbies:        true,
			WithMovieSeries:    true,
			WithTravels:        true,
			WithSports:         true,
		})
		if err != nil {
			return err
		}

		err = payload.Validate(userDetail)
		if err != nil {
			return err
		}

		err = repository.UpdateHobbies(ctx, tx, payload.Hobbies)
		if err != nil {
			return err
		}
		err = repository.UpdateMovieSeries(ctx, tx, payload.MovieSeries)
		if err != nil {
			return err
		}
		err = repository.UpdateTravelings(ctx, tx, payload.Travels)
		if err != nil {
			return err
		}
		err = repository.UpdateSports(ctx, tx, payload.Sports)
		if err != nil {
			return err
		}

		return nil
	})
}

func DeleteInterest(ctx context.Context, requestId string, payload entities.DeleteInterest) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		userDetail, err := repository.GetUserDetailById(ctx, tx, requestId, entities.GetUserDetailOption{
			PessimisticLocking: true,
			WithHobbies:        true,
			WithMovieSeries:    true,
			WithTravels:        true,
			WithSports:         true,
		})
		if err != nil {
			return err
		}

		err = payload.ValidateIds(userDetail)
		if err != nil {
			return err
		}

		err = repository.DeleteHobbiesByIds(ctx, tx, payload.HobbieIds)
		if err != nil {
			return err
		}
		err = repository.DeleteMovieSeriesByIds(ctx, tx, payload.MovieSerieIds)
		if err != nil {
			return err
		}
		err = repository.DeleteTravelingByIds(ctx, tx, payload.TravelIds)
		if err != nil {
			return err
		}
		err = repository.DeleteSportByIds(ctx, tx, payload.SportIds)
		if err != nil {
			return err
		}

		return nil

	})
}
