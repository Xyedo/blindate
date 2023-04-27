package user

import (
	"github.com/xyedo/blindate/pkg/user/dtos"
	"github.com/xyedo/blindate/pkg/user/entities"
)

type Usecase interface {
	CreateUser(dtos.RegisterUser) (string, error)
	GetUserById(dtos.GetUserDetail) (entities.User, error)
	UpdateUser(dtos.UpdateUser) error
	CreateNewProfilePic(dtos.RegisterProfilePicture) (string, error)
}
