package usecase

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func CreateUserDetail(ctx context.Context, requestId string, payload entities.CreateUserDetail) (string, error) {
	var returnedId string
	err := pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		_, err := repository.GetUserById(ctx, tx, requestId)
		if err != nil {
			return err
		}

		id, err := repository.StoreUserDetail(ctx, tx,
			entities.UserDetail{
				UserId:           requestId,
				Geog:             payload.Geog,
				Bio:              payload.Bio,
				LastOnline:       time.Now(),
				Gender:           entities.Gender(payload.Gender),
				FromLoc:          payload.FromLoc,
				Height:           payload.Height,
				EducationLevel:   payload.EducationLevel,
				Drinking:         payload.Drinking,
				Smoking:          payload.Smoking,
				RelationshipPref: payload.RelationshipPref,
				LookingFor:       payload.LookingFor,
				Zodiac:           payload.Zodiac,
				Kids:             payload.Kids,
				Work:             payload.Work,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
				Version:          1,
			},
		)
		if err != nil {
			return err
		}

		returnedId = id
		return nil
	},
	)
	if err != nil {
		return "", err
	}

	return returnedId, nil
}

func GetUserDetail(ctx context.Context, requestId, userId string) (entities.UserDetail, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return entities.UserDetail{}, err
	}

	defer conn.Release()

	//TODO: can check another userId if match/revealed
	return repository.GetUserDetailById(ctx, conn, requestId, entities.GetUserDetailOption{
		WithHobbies:     true,
		WithMovieSeries: true,
		WithTravels:     true,
		WithSports:      true,
	})

}

func UpdateUserDetailById(ctx context.Context, requestId string, payload entities.UpdateUserDetail) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		_, err := repository.GetUserDetailById(ctx, tx,
			requestId,
			entities.GetUserDetailOption{
				PessimisticLocking: true,
			},
		)
		if err != nil {
			return err
		}

		return repository.UpdateUserDetailById(ctx, tx, requestId, payload)
	})
}
