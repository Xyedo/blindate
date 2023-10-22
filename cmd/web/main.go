package main

import (
	"context"
	"log"
	"time"

	"github.com/xyedo/blindate/internal/infrastructure"
	"github.com/xyedo/blindate/internal/infrastructure/httpserver"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	infrastructure.LoadConfig(".env")

	pool, err := pg.InitPool(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer pool.Close()

	cancel()
	err = httpserver.NewEcho().Listen()
	if err != nil {
		log.Fatalln(err)
	}
}
