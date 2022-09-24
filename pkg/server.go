package pkg

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	myhttp "github.com/xyedo/blindate/pkg/http"
)

func NewServer() error {
	srv := &http.Server{
		Addr:         ":8080",
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
