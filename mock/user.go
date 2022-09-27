package mock

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct{}

func (u UserRepository) InsertUser(user *domain.User) error {
	switch {
	case user.Email == "dupli23@gmail.com":
		return &pq.Error{
			Code:    "23505",
			Message: "users_email is duplicated",
		}
	default:
		user.ID = "1"
		return nil
	}
}
func (u UserRepository) UpdateUser(user *domain.User) error {

	return nil
}
func (u UserRepository) GetUserById(id string) (*domain.User, error) {
	if id == "8c540e20-75d1-4513-a8e3-72dc4bc68619" {
		user := domain.User{
			ID:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			FullName: "Uncle Bob",
			Email:    "bob@example.com",
			Password: "pa55word",
			Active:   true,
			Dob:      time.Date(2000, time.August, 23, 0, 0, 0, 0, time.UTC),
		}
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			return nil, err
		}
		user.HashedPassword = string(hashedPass)
		user.Password = ""

		return &user, nil
	}
	return nil, sql.ErrNoRows
}
func (u UserRepository) GetUserByEmail(email string) (*domain.User, error) {

	switch email {
	case "notFound@example.com":
		return nil, sql.ErrNoRows
	default:
		user := domain.User{
			ID:       "8c540e20-75d1-4513-a8e3-72dc4bc68619",
			FullName: "Uncle Bob",
			Email:    "bob@example.com",
			Password: "pa55word",
			Active:   true,
			Dob:      time.Date(2000, time.August, 23, 0, 0, 0, 0, time.UTC),
		}
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			return nil, err
		}
		user.HashedPassword = string(hashedPass)
		user.Password = ""

		return &user, nil
	}
}
