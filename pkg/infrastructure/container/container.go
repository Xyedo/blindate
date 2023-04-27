package container

import (
	"github.com/jmoiron/sqlx"
	s3manager "github.com/xyedo/blindate/pkg/attachment/s3"
	"github.com/xyedo/blindate/pkg/infrastructure"
	"github.com/xyedo/blindate/pkg/user"
	userHttpV1 "github.com/xyedo/blindate/pkg/user/interfaces/http/v1"
	userRepo "github.com/xyedo/blindate/pkg/user/pg-repository"
	userUsecase "github.com/xyedo/blindate/pkg/user/usecase"
)

type Service struct {
	UserUC user.Usecase
}

func New(db *sqlx.DB, config infrastructure.Config) Service {
	attachment := s3manager.NewS3(config.BucketName)

	userRepo := userRepo.NewUser(db)
	userUC := userUsecase.NewUserUsecase(userRepo)
	userHandler := userHttpV1.NewUserHandler(userUC, attachment)

	return Service{
		UserUC: userUC,
	}
}
