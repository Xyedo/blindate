package entity

import (
	"database/sql"
	"time"
)

type Match struct {
	Id            string       `db:"id"`
	RequestFrom   string       `db:"request_from"`
	RequestTo     string       `db:"request_to"`
	RequestStatus string       `db:"request_status"`
	CreatedAt     time.Time    `db:"created_at"`
	AcceptedAt    sql.NullTime `db:"accepted_at"`
	RevealStatus  string       `db:"reveal_status"`
	RevealedAt    sql.NullTime `db:"revealed_at"`
}
