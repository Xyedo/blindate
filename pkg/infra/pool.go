package infra

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func (cfg *Config) OpenPgDb() (*sqlx.DB, error) {
	db := sqlx.MustOpen("postgres", cfg.DbConf.Dsn)
	db.SetMaxOpenConns(cfg.DbConf.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DbConf.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DbConf.MaxIdleTime)
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
