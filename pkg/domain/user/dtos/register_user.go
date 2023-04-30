package userDTOs

import (
	"time"
)

type RegisterUser struct {
	FullName string
	Alias    string
	Email    string
	Password string
	Dob      time.Time
}
