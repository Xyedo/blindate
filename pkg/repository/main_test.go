package repository

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var testQuery *sqlx.DB

func TestMain(m *testing.M) {
	conn, err := sqlx.Open("postgres", "postgres://blindate:pa55word@localhost:5433/blindate?sslmode=disable")

	if err != nil {
		log.Fatal("cannot connect to db:", err)
		return
	}
	testQuery = conn
	c := m.Run()
	os.Exit(c)

}
