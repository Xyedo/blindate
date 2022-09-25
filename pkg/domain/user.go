package domain

import (
	"time"
)

type User struct {
	ID        string    `db:"id" json:"Id"`
	FullName  string    `db:"full_name" json:"fullName"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password"`
	Active    bool      `db:"active" json:"-"`
	Dob       time.Time `db:"dob" json:"dob"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type UserRepository interface {
	InsertUser(user *User) error
}

type UserService interface {
	CreateUser(user *User) error
}
