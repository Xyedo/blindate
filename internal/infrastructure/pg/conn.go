package pg

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/xyedo/blindate/internal/infrastructure"
)

var (
	connection *pgx.Conn
	once       sync.Once
)

func InitConnection() {
	config := infrastructure.Config.DbConf
	conn, err := pgx.ConnectConfig(context.TODO(), &pgx.ConnConfig{
		Config: pgconn.Config{
			Host:     config.Host,
			Port:     uint16(config.Port),
			Database: config.Database,
			User:     config.User,
			Password: config.Password,
		},
	})

	if err != nil {
		panic(err)
	}

	connection = conn

}

func Get() *pgx.Conn {
	once.Do(func() {
		InitConnection()
	})
	return connection
}
