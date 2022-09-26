package domain

import (
	"time"
)

type User struct {
	ID        string    `db:"id" json:"Id"`
	FullName  string    `db:"full_name" json:"fullName"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"password,omitempty"`
	Active    bool      `db:"active" json:"-"`
	Dob       time.Time `db:"dob" json:"dob"`
	CreatedAt time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

type UserRepository interface {
	InsertUser(user *User) error
	GetUserByEmail(email string) (*User, error)
}

type UserService interface {
	CreateUser(user *User) error
	VerifyCredential(email, password string) error
}
