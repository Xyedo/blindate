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
	err := u.userRepository.InsertUser(user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			log.Println(pqErr.Code, pqErr.Message)
			if pqErr.Code == "23505" && strings.Contains(pqErr.Message, "users_email") {
				return domain.ErrDuplicateEmail
			}

			log.Panic(pqErr)
			return err
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		log.Panic(err)
		return err

	}
	return nil
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
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domain.ErrNotMatchCredential
		}
		panic(err)
	}
	return nil

}
