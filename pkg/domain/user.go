package domain

import (
	"time"
)

type User struct {
	ID             string    `db:"id" json:"Id"`
	FullName       string    `db:"full_name" json:"fullName"`
	Email          string    `db:"email" json:"email"`
	Password       string    `db:"-" json:"-"`
	HashedPassword string    `db:"password" json:"-"`
	Active         bool      `db:"active" json:"-"`
	Dob            time.Time `db:"dob" json:"dob"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

type UserRepository interface {
	InsertUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	GetUserById(id string) (*User, error)
	UpdateUser(user *User) error
}

type UserService interface {
	CreateUser(user *User) error
	VerifyCredential(email, password string) error
	GetUserById(id string) (*User, error)
	UpdateUser(user *User) error
}
