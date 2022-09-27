package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(userRepo domain.UserRepository) *User {
	return &User{
		userRepository: userRepo,
	}
}

type User struct {
	userRepository domain.UserRepository
}

func (u *User) CreateUser(user *domain.User) error {
	hashedPass, err := hashAndSalt(user.Password)
	if err != nil {
		return err
	}
	user.HashedPassword = hashedPass
	user.Password = ""

	err = u.userRepository.InsertUser(user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			log.Println(pqErr.Code, pqErr.Message)
			if pqErr.Code == "23505" && strings.Contains(pqErr.Message, "users_email") {
				return domain.ErrDuplicateEmail
			}
			return err
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	return nil
}
func (u *User) GetUserById(id string) (*domain.User, error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccesingDB
		}
		return nil, err
	}
	return user, nil
}
func (u *User) VerifyCredential(email, password string) error {
	user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotMatchCredential
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domain.ErrNotMatchCredential
		}
		return err
	}
	return nil

}
func (u *User) UpdateUser(user *domain.User) error {
	if user.Password != "" {
		hashedPass, err := hashAndSalt(user.Password)
		if err != nil {
			return err
		}
		user.HashedPassword = hashedPass
		user.Password = ""
	}
	err := u.userRepository.UpdateUser(user)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrUserNotFound
		}
		return err
	}
	return nil
}

func hashAndSalt(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPass), err
}
