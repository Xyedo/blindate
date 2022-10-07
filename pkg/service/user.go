package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type User interface {
	CreateUser(newUser *domain.User) error
	VerifyCredential(email, password string) (string, error)
	GetUserById(id string) (*domain.User, error)
	UpdateUser(user *domain.User) error
}

func NewUser(userRepo repository.User) *user {
	return &user{
		userRepository: userRepo,
	}
}

type user struct {
	userRepository repository.User
}

func (u *user) CreateUser(newUser *domain.User) error {
	hashedPass, err := hashAndSalt(newUser.Password)
	if err != nil {
		return err
	}
	newUser.HashedPassword = hashedPass
	newUser.Password = ""

	err = u.userRepository.InsertUser(newUser)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return domain.ErrUniqueConstraint23505
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
func (u *user) GetUserById(id string) (*domain.User, error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccesingDB
		}
		return nil, err
	}
	return user, nil
}
func (u *user) VerifyCredential(email, password string) (string, error) {
	user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrNotMatchCredential
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", domain.ErrNotMatchCredential
		}
		return "", err
	}
	return user.ID, nil

}
func (u *user) UpdateUser(user *domain.User) error {
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
			return domain.ErrResourceNotFound
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
