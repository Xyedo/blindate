package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
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
			if pqErr.Code == "42701" && strings.Contains(pqErr.Message, "email") {
				return domain.ErrDuplicateEmail
			}
			//unhandled postgre error
			//TODO REMOVE THIS
			log.Panic(pqErr)
			return err
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		//unhandled unknown error
		//TODO REMOVE THIS
		log.Panic(err)
		return err

	}
	return nil
}
