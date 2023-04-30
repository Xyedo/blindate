package user

import (
	userDTOs "github.com/xyedo/blindate/pkg/domain/user/dtos"
	userEntities "github.com/xyedo/blindate/pkg/domain/user/entities"
)

type Usecase interface {
	CreateUser(userDTOs.RegisterUser) (string, error)
	GetUserById(userDTOs.GetUserDetail) (userEntities.User, error)
	UpdateUser(userDTOs.UpdateUser) error
	CreateNewProfilePic(userDTOs.RegisterProfilePicture) (string, error)
}
