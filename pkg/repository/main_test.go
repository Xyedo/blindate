package repository_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQuery *sqlx.DB

func TestMain(m *testing.M) {
	godotenv.Load("../../.env.dev")
	conn, err := sqlx.Open("postgres", os.Getenv("POSTGRE_DB_DSN_TEST"))

	if err != nil {
		log.Fatal("cannot connect to db:", err)
		return
	}
	testQuery = conn
	c := m.Run()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	testQuery.MustExecContext(ctx, `DELETE FROM users WHERE 1=1`)
	testQuery.MustExecContext(ctx, `DELETE FROM authentications WHERE 1=1`)
	os.Exit(c)

}
