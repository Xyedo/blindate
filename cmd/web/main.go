package main

import (
	"log"

	"github.com/xyedo/blindate/internal/infrastructure"
	"github.com/xyedo/blindate/internal/infrastructure/httpserver"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func main() {
	infrastructure.LoadConfig(".env")
	pg.InitConnection()

	server := httpserver.NewEcho()

	err := server.Listen()
	if err != nil {
		log.Fatal(err)
	}
}
