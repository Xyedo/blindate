package postgre

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/infrastructure"
)

func OpenDB(cfg infrastructure.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbConf.Host,
		cfg.DbConf.Port,
		cfg.DbConf.User,
		cfg.DbConf.Password,
		cfg.DbConf.Database,
	)
	
	db := sqlx.MustOpen("postgres", dsn)
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
