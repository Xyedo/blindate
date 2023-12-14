package usecase

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/common/ids"
	attachmentEntities "github.com/xyedo/blindate/internal/domain/attachment/entities"
	attachmentRepo "github.com/xyedo/blindate/internal/domain/attachment/repository"
	"github.com/xyedo/blindate/internal/domain/attachment/s3"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	userentities "github.com/xyedo/blindate/internal/domain/user/entities"
	userRepo "github.com/xyedo/blindate/internal/domain/user/repository"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func AddPhoto(ctx context.Context, requestId string, header *multipart.FileHeader) (string, error) {
	photo, err := header.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = photo.Close() }()

	contentType, err := userentities.SanitizeMimeType(photo, []string{"image/jpeg", "image/png", "image/webp"})
	if err != nil {
		return "", err
	}

	var photoId string
	timeNow := time.Now()
	err = pg.Transaction(ctx, pgx.TxOptions{}, func(tx pg.Querier) error {
		userDetail, err := userRepo.GetUserDetailById(ctx, tx, requestId, userentities.GetUserDetailOption{
			PessimisticLocking:  true,
			WithProfilePictures: true,
		})
		if err != nil {
			return err
		}

		if len(userDetail.ProfilePictures) > 5 {
			return apperror.BadPayloadWithPayloadMap(apperror.PayloadMap{
				Payloads: []apperror.ErrorPayload{
					{
						Code: entities.PhotoTooMuch,
						Details: map[string][]string{
							"file": {"exceeding profile-photo upload"},
						},
					},
				},
			})
		}
		objectKey, err := s3.Manager.UploadAttachment(ctx, photo, attachmentEntities.Attachment{
			File:        photo,
			ContentType: contentType,
			Prefix:      "/user/" + requestId + "/photos",
			Ext:         filepath.Ext(header.Filename),
		})
		if err != nil {
			return err
		}

		fileId, err := attachmentRepo.InsertFile(ctx, tx, attachmentEntities.File{
			Id:        ids.File(),
			FileType:  attachmentEntities.FileTypePhotoProfile,
			BlobLink:  objectKey,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Version:   1,
		})
		if err != nil {
			return err
		}

		err = userRepo.UpdateProfilePictureSelectedToFalseByUserId(ctx, tx, userDetail.UserId)
		if err != nil {
			return err
		}

		returnedProfilePictureId, err := userRepo.InsertProfilePicture(ctx, tx, userentities.ProfilePicture{
			Id:       ids.ProfilePicture(),
			UserId:   userDetail.UserId,
			Selected: true,
			FileId:   fileId,
		})
		if err != nil {
			return err
		}
		photoId = returnedProfilePictureId
		return nil

	})
	if err != nil {
		return "", err
	}

	return photoId, nil
}
