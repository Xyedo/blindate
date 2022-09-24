package main

import (
	"flag"
	"os"

	"github.com/xyedo/blindate/pkg"
)

func main() {
	var cfg pkg.Config

	flag.IntVar(&cfg.Port, "port", 8080, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environtment (development | staging | production)")

	flag.StringVar(&cfg.Db.Dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PgSQL dsn")
	flag.IntVar(&cfg.Db.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.Db.MaxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.Db.MaxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	pkg.NewServer(cfg)
}
