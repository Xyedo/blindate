package onlineEntities

import "time"

type Online struct {
	UserId     string    `db:"user_id"`
	LastOnline time.Time `db:"last_online"`
	IsOnline   bool      `db:"is_online"`
}
