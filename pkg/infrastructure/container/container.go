package container

import (
	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/internal/security"
	"github.com/xyedo/blindate/pkg/domain/attachment"
	s3manager "github.com/xyedo/blindate/pkg/domain/attachment/s3"
	"github.com/xyedo/blindate/pkg/domain/authentication"
	authRepo "github.com/xyedo/blindate/pkg/domain/authentication/pg-repository"
	authUsecase "github.com/xyedo/blindate/pkg/domain/authentication/usecase"
	basicinfo "github.com/xyedo/blindate/pkg/domain/basic-info"
	basicinfoRepo "github.com/xyedo/blindate/pkg/domain/basic-info/pg-repository"
	basicinfoUsecase "github.com/xyedo/blindate/pkg/domain/basic-info/usecase"
	"github.com/xyedo/blindate/pkg/domain/gateway"
	gatewayUsecase "github.com/xyedo/blindate/pkg/domain/gateway/usecase"
	"github.com/xyedo/blindate/pkg/domain/user"
	userRepo "github.com/xyedo/blindate/pkg/domain/user/pg-repository"
	userUsecase "github.com/xyedo/blindate/pkg/domain/user/usecase"
	"github.com/xyedo/blindate/pkg/infrastructure"
)

type Container struct {
	AttachmentManager attachment.Repository
	Jwt               *security.Jwt
	GatewaySession    gateway.Session
	UserUC            user.Usecase
	AuthUC            authentication.Usecase
	BasicInfoUC       basicinfo.Usecase
}

func New(db *sqlx.DB, config infrastructure.Config) Container {
	attachment := s3manager.NewS3(config.BucketName)

	jwt := security.NewJwt(
		config.Token.AccessSecret,
		config.Token.RefreshSecret,
		config.Token.AccessExpires,
		config.Token.RefreshExpires,
	)

	userRepo := userRepo.New(db)
	authRepo := authRepo.New(db)
	basicInfoRepo := basicinfoRepo.New(db)

	authUC := authUsecase.New(authRepo, userRepo, jwt)
	userUC := userUsecase.New(userRepo)
	basicinfoUC := basicinfoUsecase.New(basicInfoRepo)

	gatewaySession := gatewayUsecase.NewSession()
	return Container{
		AttachmentManager: attachment,
		Jwt:               jwt,
		GatewaySession:    gatewaySession,
		UserUC:            userUC,
		AuthUC:            authUC,
		BasicInfoUC:       basicinfoUC,
	}
}
