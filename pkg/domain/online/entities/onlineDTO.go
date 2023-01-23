package onlineEntity

import "time"

// Online one to one with user
type DTO struct {
	UserId     string    `json:"-" db:"user_id"`
	LastOnline time.Time `json:"lastOnline" db:"last_online"`
	IsOnline   bool      `json:"isOnline" db:"is_online"`
}
