package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	authHttpV1 "github.com/xyedo/blindate/pkg/domain/authentication/interfaces/http/v1"
	basicInfoHttpV1 "github.com/xyedo/blindate/pkg/domain/basic-info/interfaces/http/v1"
	gatewayHttpV1 "github.com/xyedo/blindate/pkg/domain/gateway/interfaces/http/v1"
	interestHttpV1 "github.com/xyedo/blindate/pkg/domain/interest/interfaces/http/v1"
	onlineHttpV1 "github.com/xyedo/blindate/pkg/domain/online/interfaces/http/v1"
	userHttpV1 "github.com/xyedo/blindate/pkg/domain/user/interfaces/http/v1"
	"github.com/xyedo/blindate/pkg/infrastructure"
	"github.com/xyedo/blindate/pkg/infrastructure/container"
)

func NewServer(config infrastructure.Config, container *container.Container) httpServer {
	gin.EnableJsonDecoderDisallowUnknownFields()

	ginEngine := gin.New()
	ginEngine.HandleMethodNotAllowed = true
	ginEngine.MaxMultipartMemory = 8 << 20

	return httpServer{
		config:    config,
		gin:       ginEngine,
		container: container,
	}
}

type httpServer struct {
	config    infrastructure.Config
	gin       *gin.Engine
	container *container.Container
}

func (s *httpServer) Listen() error {
	s.gin.NoMethod(noMethod)
	s.gin.NoRoute(noFound)
	if s.config.Env == "development" {
		s.gin.Use(gin.Logger())
	}
	s.gin.Use(gin.Recovery())

	s.handlerV1()

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:      s.gin,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownErr := make(chan error)
	go gracefulShutDown(shutdownErr, srv)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	defer close(shutdownErr)
	return <-shutdownErr

}

func (h *httpServer) handlerV1() {
	v1 := h.gin.Group("/api/v1")
	v1.Use(cors.Default())

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	userHandler := userHttpV1.New(h.config, h.container.UserUC, h.container.AttachmentManager)
	authHandler := authHttpV1.New(h.config, h.container.AuthUC)
	gatewayHandler := gatewayHttpV1.New(h.config, h.container.GatewaySession)
	basicInfoHandler := basicInfoHttpV1.New(h.config, h.container.BasicInfoUC)
	onlineHandler := onlineHttpV1.New(h.config, h.container.OnlineUC)
	interestHandler := interestHttpV1.New(h.config, h.container.InterestUC)

	go gatewayHandler.Listen()

	gatewayHandler.Handler(v1, h.container.Jwt)
	userHandler.Handler(v1, h.container.Jwt)
	authHandler.Handler(v1)
	basicInfoHandler.Handler(v1, h.container.Jwt)
	onlineHandler.Handler(v1, h.container.Jwt)
	interestHandler.Handler(v1, h.container.Jwt)

}

func gracefulShutDown(shutdownError chan<- error, server *http.Server) {
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownError <- server.Shutdown(ctx)
}

func noFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  "failed",
		"message": "not found",
	})
}

func noMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{
		"status":  "failed",
		"message": "method not allowed",
	})
}
