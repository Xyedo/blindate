package pg

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/xyedo/blindate/internal/infrastructure"
)

var (
	connection *pgx.Conn
	once       sync.Once
)

func InitConnection() {
	once.Do(func() {
		initConnection()
	})
}

func Get() *pgx.Conn {
	once.Do(func() {
		initConnection()
	})
	return connection
}

func initConnection() {
	config := infrastructure.Config.DbConf

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	conn, err := pgx.Connect(context.TODO(), connStr)

	if err != nil {
		panic(err)
	}

	connection = conn

}
