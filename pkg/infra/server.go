package infra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xyedo/blindate/pkg/api"
)

type Config struct {
	Port       int
	Env        string
	BucketName string
	DbConf     struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	Token struct {
		AccessSecret   string
		RefreshSecret  string
		AccessExpires  string
		RefreshExpires string
	}
}

func (cfg *Config) NewServer(route api.Route) error {
	handler := api.Routes(route)
	srv := &http.Server{

		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler:      handler,
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

func gracefulShutDown(shutdownError chan<- error, server *http.Server) {
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownError <- server.Shutdown(ctx)
}
