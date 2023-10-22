package usecase

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func CreateBasicInfo(ctx context.Context, requestId string, payload entities.CreateBasicInfo) (string, error) {
	var returnedId string
	err := pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		returnedUser, err := repository.GetUserById(ctx, tx, requestId)
		if err != nil {
			return err
		}

		id, err := repository.StoreBasicInfo(ctx, tx,
			entities.BasicInfo{
				UserId:           returnedUser.Id,
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

func GetBasicInfoById(ctx context.Context, requestId, userId string) (entities.BasicInfo, error) {
	//TODO: can check another userId if match/revealed
	return repository.GetBasicInfoById(ctx, pg.GetPool(ctx), requestId)

}

func UpdateBasicInfoById(ctx context.Context, requestId string, payload entities.UpdateBasicInfo) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		_, err := repository.GetBasicInfoById(ctx, tx,
			requestId,
			entities.GetBasicInfoOption{
				PessimisticLocking: true,
			},
		)
		if err != nil {
			return err
		}

		return repository.UpdateBasicInfoById(ctx, tx, requestId, payload)
	})
}
