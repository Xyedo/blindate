package locationEntity

import "time"

type DAO struct {
	UserId    string    `db:"user_id"`
	Geog      string    `db:"geog"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
