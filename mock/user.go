package mock

import (
	"github.com/xyedo/blindate/pkg/domain"
)

type UserService struct{}

func (u UserService) CreateUser(user *domain.User) error {
	switch {
	case user.Email == "dupli@example.com":
		return domain.ErrDuplicateEmail
	case len(user.Password) > 100:
		return domain.ErrTooLongAccesingDB
	}
	return nil

}
