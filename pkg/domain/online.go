package domain

import "time"

type Online struct {
	UserId     string    `json:"-" db:"user_id"`
	LastOnline time.Time `json:"lastOnline" db:"last_online"`
	IsOnline   bool      `json:"isOnline" db:"is_online"`
}
