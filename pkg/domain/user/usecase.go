package user

import (
	userDTOs "github.com/xyedo/blindate/pkg/domain/user/dtos"
)

type Usecase interface {
	CreateUser(userDTOs.RegisterUser) (string, error)
	GetUserById(userDTOs.GetUserDetail) (userDTOs.User, error)
	UpdateUser(userDTOs.UpdateUser) error
	CreateNewProfilePic(userDTOs.RegisterProfilePicture) (string, error)
}
