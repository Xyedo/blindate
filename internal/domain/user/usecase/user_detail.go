package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	attachmentRepo "github.com/xyedo/blindate/internal/domain/attachment/repository"
	"github.com/xyedo/blindate/internal/domain/attachment/s3"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	userRepo "github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func CreateUserDetail(ctx context.Context, requestId string, payload entities.CreateUserDetail) (string, error) {
	var returnedId string
	err := pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		_, err := userRepo.GetUserById(ctx, tx, requestId)
		if err != nil {
			return err
		}

		id, err := userRepo.StoreUserDetail(ctx, tx,
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

func GetUserDetail(ctx context.Context, requestId, userId string) (entities.UserDetail, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return entities.UserDetail{}, err
	}

	defer conn.Release()

	//TODO: can check another userId if match/revealed
	userDetail, err := userRepo.GetUserDetailById(ctx, conn, requestId, entities.GetUserDetailOption{
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
		files, err := attachmentRepo.GetFileByIds(ctx, conn, fileIds)
		if err != nil {
			return entities.UserDetail{}, err
		}

		var wg sync.WaitGroup
		errs := make([]error, len(files))

		wg.Add(len(files))
		for i := 0; i < len(files); i++ {
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				presignedURL, err := s3.Manager.GetPresignedUrl(ctx, files[i].BlobLink, 1*time.Hour)
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

func UpdateUserDetailById(ctx context.Context, requestId string, payload entities.UpdateUserDetail) error {
	return pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		_, err := userRepo.GetUserDetailById(ctx, tx,
			requestId,
			entities.GetUserDetailOption{
				PessimisticLocking: true,
			},
		)
		if err != nil {
			return err
		}

		return userRepo.UpdateUserDetailById(ctx, tx, requestId, payload)
	})
}

// GetUserDetails
// its For External Domain want to Get User Details
func GetUserDetails(ctx context.Context, userIds []string) (entities.UserDetails, error) {
	conn, err := pg.GetConnectionPool(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Release()

	userDetails, err := userRepo.FindUserDetailByIds(ctx, conn, userIds)
	if err != nil {
		return nil, err
	}

	fileIds, fileIdToIdx := userDetails.ToFileIds()
	if len(fileIds) > 0 {
		files, err := attachmentRepo.GetFileByIds(ctx, conn, fileIds)
		if err != nil {
			return nil, err
		}

		var wg sync.WaitGroup
		errs := make([]error, len(files))

		wg.Add(len(files))
		for i := 0; i < len(files); i++ {
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				presignedURL, err := s3.Manager.GetPresignedUrl(ctx, files[i].BlobLink, 1*time.Hour)
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
