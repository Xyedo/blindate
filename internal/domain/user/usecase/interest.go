package usecase

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func (uc *User) CreateInterest(ctx context.Context, requestId string, payload entities.CreateInterest) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		userDetail, err := uc.repo.GetUserDetailById(ctx, tx, requestId, entities.GetUserDetailOption{
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

		err = uc.repo.StoreHobbiesByUserId(ctx, tx, requestId, payload.ToHobbies(requestId))
		if err != nil {
			return err
		}

		err = uc.repo.StoreMovieSeriesByUserId(ctx, tx, requestId, payload.ToMovieSeries(requestId))
		if err != nil {
			return err
		}

		err = uc.repo.StoreTravelingsByUserId(ctx, tx, requestId, payload.ToTravels(requestId))
		if err != nil {
			return err
		}

		err = uc.repo.StoreSportsByUserId(ctx, tx, requestId, payload.ToSports(requestId))
		if err != nil {
			return err
		}

		return nil
	})

}

func (uc *User) UpdateInterest(ctx context.Context, requestId string, payload entities.UpdateInterest) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		userDetail, err := uc.repo.GetUserDetailById(ctx, tx, requestId, entities.GetUserDetailOption{
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

		err = uc.repo.UpdateHobbies(ctx, tx, payload.Hobbies)
		if err != nil {
			return err
		}
		err = uc.repo.UpdateMovieSeries(ctx, tx, payload.MovieSeries)
		if err != nil {
			return err
		}
		err = uc.repo.UpdateTravelings(ctx, tx, payload.Travels)
		if err != nil {
			return err
		}
		err = uc.repo.UpdateSports(ctx, tx, payload.Sports)
		if err != nil {
			return err
		}

		return nil
	})
}

func (uc *User) DeleteInterest(ctx context.Context, requestId string, payload entities.DeleteInterest) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		userDetail, err := uc.repo.GetUserDetailById(ctx, tx, requestId, entities.GetUserDetailOption{
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

		err = uc.repo.DeleteHobbiesByIds(ctx, tx, payload.HobbieIds)
		if err != nil {
			return err
		}
		err = uc.repo.DeleteMovieSeriesByIds(ctx, tx, payload.MovieSerieIds)
		if err != nil {
			return err
		}
		err = uc.repo.DeleteTravelingByIds(ctx, tx, payload.TravelIds)
		if err != nil {
			return err
		}
		err = uc.repo.DeleteSportByIds(ctx, tx, payload.SportIds)
		if err != nil {
			return err
		}

		return nil

	})
}
