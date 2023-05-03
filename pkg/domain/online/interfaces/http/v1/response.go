package v1

import "time"

type getUserOnlineResponse struct {
	UserId     string    `json:"user_id"`
	LastOnline time.Time `json:"last_online"`
	IsOnline   bool      `json:"is_online"`
}
