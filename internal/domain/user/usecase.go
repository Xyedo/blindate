package user

import (
	"context"
	"mime/multipart"

	"github.com/xyedo/blindate/internal/domain/user/entities"
)

type Usecase interface {
	RegisterUser(ctx context.Context, id string) error
	DeleteUser(ctx context.Context, id string) error

	CreateUserDetail(ctx context.Context, requestId string, payload entities.CreateUserDetail) (string, error)
	GetUserDetail(ctx context.Context, requestId string, userId string) (entities.UserDetail, error)
	UpdateUserDetailById(ctx context.Context, requestId string, payload entities.UpdateUserDetail) error

	CreateInterest(ctx context.Context, requestId string, payload entities.CreateInterest) error
	UpdateInterest(ctx context.Context, requestId string, payload entities.UpdateInterest) error
	DeleteInterest(ctx context.Context, requestId string, payload entities.DeleteInterest) error

	AddPhoto(ctx context.Context, requestId string, header *multipart.FileHeader) (string, error)
}
