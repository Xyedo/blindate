package httpserver

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	repositoryAttachment "github.com/xyedo/blindate/internal/domain/attachment/repository"
	"github.com/xyedo/blindate/internal/domain/attachment/storage"
	usecaseAttachment "github.com/xyedo/blindate/internal/domain/attachment/usecase"
	conversationHandler "github.com/xyedo/blindate/internal/domain/conversation/handler"
	repositoryConversation "github.com/xyedo/blindate/internal/domain/conversation/repository"
	usecaseConversation "github.com/xyedo/blindate/internal/domain/conversation/usecase"
	matchHandler "github.com/xyedo/blindate/internal/domain/match/handler"
	repositoryMatch "github.com/xyedo/blindate/internal/domain/match/repository"
	usecaseMatch "github.com/xyedo/blindate/internal/domain/match/usecase"
	externalUserHandler "github.com/xyedo/blindate/internal/domain/user/handler/external"
	internalUserHandler "github.com/xyedo/blindate/internal/domain/user/handler/intrnl"
	repositoryUser "github.com/xyedo/blindate/internal/domain/user/repository"
	usecaseUser "github.com/xyedo/blindate/internal/domain/user/usecase"

	"github.com/xyedo/blindate/internal/infrastructure"
	echomiddleware "github.com/xyedo/blindate/internal/infrastructure/httpserver/echo-middleware"
)

func NewEcho() *Server {
	e := echo.New()
	e.HTTPErrorHandler = EchoErrorHandler

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
		middleware.CORS(),
		middleware.BodyLimit("4M"),
		middleware.ContextTimeout(3*time.Second))

	e.GET("/healthcheck", func(c echo.Context) error {
		return nil
	})

	//repository
	userRepo := repositoryUser.User{}
	matchRepo := repositoryMatch.Match{}
	conversationRepo := repositoryConversation.Conversation{}
	fileRepo := repositoryAttachment.File{}
	storage := storage.NewS3()

	//usecase
	attachmentUsecase := usecaseAttachment.New(fileRepo, storage)
	userUsecase := usecaseUser.New(userRepo, attachmentUsecase)
	matchUsecase := usecaseMatch.New(matchRepo, userUsecase)
	conversationUsecase := usecaseConversation.New(
		conversationRepo,
		userUsecase,
		attachmentUsecase,
	)
	//handler
	apiv1 := e.Group("/v1")
	{

		internal(apiv1, userUsecase)

		apiv1.Use(echomiddleware.Guard)

		externalUserHandler.
			New(userUsecase).
			Route(apiv1)

		matchHandler.
			New(matchUsecase).
			Route(apiv1)

		conversationHandler.
			New(conversationUsecase).
			Route(apiv1)
	}

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", infrastructure.Config.Host, infrastructure.Config.Port),
			Handler: e,
		},
	}
}

func internal(e *echo.Group, userUsecase *usecaseUser.User) {
	internal := e.Group("/internal")
	internalUserHandler.
		New(userUsecase).
		Route(internal)
}
