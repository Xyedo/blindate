package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/infrastructure"
	"github.com/xyedo/blindate/pkg/infrastructure/container"
	"github.com/xyedo/blindate/pkg/infrastructure/http"
	"github.com/xyedo/blindate/pkg/infrastructure/postgre"
)

// @title Blindate
// @version 1.0
// @description Blindate API Docs
// @contact.name hafid mahdi
// @contact.url xyedo.dev
// @contact.email hafidmahdi23@gmail.com
// @Schemes https
// @host api-dev.mceasy.com
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	var cfg infrastructure.Config
	cfg.LoadConfig(".env.dev")

	db, err := postgre.OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}(db)

	container := container.New(db, cfg)
	server := http.NewServer(cfg, container)

	err = server.Listen()
	if err != nil {
		log.Fatal(err)
	}

}
