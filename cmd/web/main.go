package main

import (
	"log"

	"github.com/xyedo/blindate/internal/infrastructure"
	"github.com/xyedo/blindate/internal/infrastructure/httpserver"
)

func main() {
	infrastructure.LoadConfig(".env")

	err := httpserver.NewEcho().Listen()
	if err != nil {
		log.Fatal(err)
	}
}
