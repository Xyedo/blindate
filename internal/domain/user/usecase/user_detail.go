package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func (uc *User) CreateUserDetail(ctx context.Context, requestId string, payload entities.CreateUserDetail) (string, error) {
	var returnedId string
	err := pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		_, err := uc.repo.GetUserById(ctx, tx, requestId)
		if err != nil {
			return err
		}

		id, err := uc.repo.StoreUserDetail(ctx, tx,
			entities.UserDetail{
				UserId:           requestId,
				Alias:            payload.Alias,
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

func (uc *User) GetUserDetail(ctx context.Context, requestId, userId string) (entities.UserDetail, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return entities.UserDetail{}, err
	}

	defer conn.Release()

	return uc.getUserDetail(ctx, conn, requestId, userId)
}

func (uc *User) UpdateUserDetailById(ctx context.Context, requestId string, payload entities.UpdateUserDetail) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		_, err := uc.repo.GetUserDetailById(ctx, tx,
			requestId,
			entities.GetUserDetailOption{
				PessimisticLocking: true,
			},
		)
		if err != nil {
			return err
		}

		return uc.repo.UpdateUserDetailById(ctx, tx, requestId, payload)
	})
}

// GetUserDetails
func (uc *User) GetUserDetails(ctx context.Context, conn pg.Querier, userIds []string) (entities.UserDetails, error) {
	userDetails, err := uc.repo.FindUserDetailByIds(ctx, conn, userIds)
	if err != nil {
		return nil, err
	}

	fileIds, fileIdToIdx := userDetails.ToFileIds()
	if len(fileIds) > 0 {
		files, err := uc.attachmentUsecase.FindFilesByIds(ctx, conn, fileIds)
		if err != nil {
			return nil, err
		}

		var wg sync.WaitGroup
		errs := make([]error, len(files))

		wg.Add(len(files))
		for i := 0; i < len(files); i++ {
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				presignedURL, err := uc.attachmentUsecase.GetAttachmentURL(ctx, files[i].BlobLink, 1*time.Hour)
				if err != nil {
					errs[i] = err
					return
				}

				if idxs, ok := fileIdToIdx[files[i].Id]; ok {
					userDetails[idxs[0]].ProfilePictures[idxs[1]].SetPresignedURL(presignedURL)
				}
			}(i, &wg)
		}
		wg.Wait()

		for _, err := range errs {
			if err != nil {
				return nil, err
			}
		}
	}

	return userDetails, nil
}

func (uc *User) GetUserDetailByUserId(ctx context.Context, conn pg.Querier, userId string) (entities.UserDetail, error) {
	return uc.getUserDetail(ctx, conn, userId, userId)
}

func (uc *User) getUserDetail(ctx context.Context, conn pg.Querier, requestId, userId string) (entities.UserDetail, error) {
	//TODO: can check another userId if match/revealed
	userDetail, err := uc.repo.GetUserDetailById(ctx, conn, requestId, entities.GetUserDetailOption{
		WithHobbies:         true,
		WithMovieSeries:     true,
		WithTravels:         true,
		WithSports:          true,
		WithProfilePictures: true,
	})
	if err != nil {
		return entities.UserDetail{}, err
	}

	if len(userDetail.ProfilePictures) > 0 {
		fileIds, fileIdToIdx := userDetail.ToFileIds()
		files, err := uc.attachmentUsecase.FindFilesByIds(ctx, conn, fileIds)
		if err != nil {
			return entities.UserDetail{}, err
		}

		var wg sync.WaitGroup
		errs := make([]error, len(files))

		wg.Add(len(files))
		for i := 0; i < len(files); i++ {
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				presignedURL, err := uc.attachmentUsecase.GetAttachmentURL(ctx, files[i].BlobLink, 1*time.Hour)
				if err != nil {
					errs[i] = err
					return
				}

				if idx, ok := fileIdToIdx[files[i].Id]; ok {
					userDetail.ProfilePictures[idx].SetPresignedURL(presignedURL)
				}
			}(i, &wg)
		}
		wg.Wait()

		for _, err := range errs {
			if err != nil {
				return entities.UserDetail{}, err
			}
		}
	}

	return userDetail, nil
}
