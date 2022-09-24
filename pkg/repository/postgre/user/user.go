package user

import "github.com/jmoiron/sqlx"

type userCon struct {
	db *sqlx.DB
}

func New() (userCon, error) {
	sqlx.Connect("postgres")
}
