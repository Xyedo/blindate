package domain

import (
	"time"
)

type User struct {
	ID             string    `db:"id" json:"id"`
	FullName       string    `db:"full_name" json:"fullName"`
	Email          string    `db:"email" json:"email"`
	Password       string    `db:"-" json:"-"`
	HashedPassword string    `db:"password" json:"-"`
	Active         bool      `db:"active" json:"-"`
	Dob            time.Time `db:"dob" json:"dob"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}
