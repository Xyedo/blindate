package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xyedo/blindate/pkg"
	"github.com/xyedo/blindate/pkg/api"
	"github.com/xyedo/blindate/pkg/repository"
	"github.com/xyedo/blindate/pkg/service"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environtment (development | staging | production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("POSTGRE_DB_DSN"), "PgSQL dsn")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := repository.NewUser(db)
	userSvc := service.NewUser(userRepo)
	userHandler := api.NewUser(userSvc)

	healthcheckHander := api.NewHealthCheck()

	route := api.Route{
		User:        userHandler,
		Healthcheck: healthcheckHander,
	}

	h := api.Routes(route)

	err = pkg.NewServer(cfg.port, h)
	if err != nil {
		log.Fatal(err)
	}

}

func openDB(cfg config) (*sqlx.DB, error) {
	db := sqlx.MustOpen("postgres", cfg.db.dsn)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, err
}
