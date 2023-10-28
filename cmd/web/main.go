package main

import (
	"context"
	"log"

	"github.com/xyedo/blindate/internal/infrastructure"
	"github.com/xyedo/blindate/internal/infrastructure/httpserver"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func main() {
	infrastructure.LoadConfig(".env")

	pool, err := pg.InitPool(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer pool.Close()

	err = httpserver.NewEcho().Listen()
	if err != nil {
		log.Fatalln(err)
	}
}
