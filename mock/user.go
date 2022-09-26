package mock

import (
	"github.com/xyedo/blindate/pkg/domain"
)

type UserService struct{}

func (u UserService) CreateUser(user *domain.User) error {
	switch {
	case user.Email == "dupli23@gmail.com":
		return domain.ErrDuplicateEmail
	case len(user.Password) > 15:
		return domain.ErrTooLongAccesingDB
	}
	user.ID = "1"
	return nil

}

func (u UserService) VerifyCredential(email, password string) error {
	switch email {
	case "notFound@example.com":
		return domain.ErrNotMatchCredential
	default:
		switch password {
		case "pa55word":
			return nil
		default:
			return domain.ErrNotMatchCredential
		}
	}

}
