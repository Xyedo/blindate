package locationEntities

import "time"

type Location struct {
	UserId    string    `db:"user_id"`
	Geog      string    `db:"geog"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
