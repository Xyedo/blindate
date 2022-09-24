package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	myhttp "github.com/xyedo/blindate/pkg/http"
)

func NewServer(cfg Config) error {
	srv := &http.Server{
		//todo add this to the config
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler:      myhttp.Routes(),
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
	return <-shutdownErr

}

func gracefulShutDown(shutdownError chan<- error, server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shutdownError <- server.Shutdown(ctx)
}
