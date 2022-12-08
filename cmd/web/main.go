package main

import (
	"flag"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/infra"
)

func main() {
	var cfg infra.Config

	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environtment (development | staging | production)")

	flag.StringVar(&cfg.DbConf.Dsn, "db-dsn", os.Getenv("POSTGRE_DB_DSN"), "PgSQL dsn")
	flag.IntVar(&cfg.DbConf.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DbConf.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DbConf.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.Token.AccessSecret, "jwt-access-secret", os.Getenv("JWT_ACCESS_SECRET_KEY"), "Jwt Access")
	flag.StringVar(&cfg.Token.RefreshSecret, "jwt-refresh-secret", os.Getenv("JWT_REFRESH_SECRET_KEY"), "Jwt Access")
	flag.StringVar(&cfg.Token.AccessExpires, "jwt-access-expires", os.Getenv("JWT_ACCESS_EXPIRES"), "Jwt Access")
	flag.StringVar(&cfg.Token.RefreshExpires, "jwt-refresh-expires", os.Getenv("JWT_REFRESH_EXPIRES"), "Jwt Access")

	flag.Parse()

	db, err := cfg.OpenPgDb()
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

	log.Fatal(cfg.NewServer(cfg.Container(db)))

}
