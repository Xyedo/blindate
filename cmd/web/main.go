package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/infrastructure"
	"github.com/xyedo/blindate/pkg/infrastructure/postgre"
)

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

	routes := cfg.Container(db)

	go routes.ListenToWsChan()
	err = cfg.NewServer(routes)
	if err != nil {
		log.Fatal(err)
	}

}
