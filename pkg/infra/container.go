package infra

import (
	"github.com/jmoiron/sqlx"
	"github.com/xyedo/blindate/pkg/api"
	"github.com/xyedo/blindate/pkg/repository"
	"github.com/xyedo/blindate/pkg/service"
)

func (cfg *Config) Container(db *sqlx.DB) api.Route {
	attachmentSvc := service.NewS3(cfg.BucketName)

	userRepo := repository.NewUser(db)
	userSvc := service.NewUser(userRepo)
	userHandler := api.NewUser(userSvc, attachmentSvc)

	healthcheckHander := api.NewHealthCheck()

	basicInfoRepo := repository.NewBasicInfo(db)
	basicInfoSvc := service.NewBasicInfo(basicInfoRepo)
	basicInfoHandler := api.NewBasicInfo(basicInfoSvc)

	locationRepo := repository.NewLocation(db)
	locationService := service.NewLocation(locationRepo)
	locationHandler := api.NewLocation(locationService)

	interestRepo := repository.NewInterest(db)
	interestSvc := service.NewInterest(interestRepo)
	interestHandler := api.NewInterest(interestSvc)

	onlineRepo := repository.NewOnline(db)
	onlineSvc := service.NewOnline(onlineRepo)
	onlineHandler := api.NewOnline(onlineSvc)

	tokenSvc := service.NewJwt(cfg.Token.AccessSecret, cfg.Token.RefreshSecret, cfg.Token.AccessExpires, cfg.Token.RefreshExpires)

	authRepo := repository.NewAuth(db)
	authSvc := service.NewAuth(authRepo, userRepo, tokenSvc)
	authHandler := api.NewAuth(authSvc)

	matchRepo := repository.NewMatch(db)
	matchSvc := service.NewMatch(matchRepo, locationRepo)
	matchHandler := api.NewMatch(matchSvc)

	convRepo := repository.NewConversation(db)
	convSvc := service.NewConversation(convRepo, matchRepo)
	convHandler := api.NewConvo(convSvc)

	chatRepp := repository.NewChat(db)
	chatSvc := service.NewChat(chatRepp, matchRepo)
	chatHandler := api.NewChat(chatSvc, attachmentSvc)

	return api.Route{
		User:           userHandler,
		Healthcheck:    healthcheckHander,
		BasicInfo:      basicInfoHandler,
		Location:       locationHandler,
		Authentication: authHandler,
		Tokenizer:      tokenSvc,
		Interest:       interestHandler,
		Online:         onlineHandler,
		Convo:          convHandler,
		Chat:           chatHandler,
		Match:          matchHandler,
	}
}
