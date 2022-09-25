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
