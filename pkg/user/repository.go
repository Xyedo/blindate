package user

import (
	"github.com/xyedo/blindate/pkg/user/dtos"
	"github.com/xyedo/blindate/pkg/user/entities"
)

type Repository interface {
	InsertUser(user dtos.RegisterUser) (string, error)
	GetUserById(id string) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	UpdateUser(user entities.User) error
	CreateProfilePicture(userId, pictureRef string, selected bool) (string, error)
	SelectProfilePicture(userId string, params dtos.ProfilePictureQuery) ([]entities.UserProfilePic, error)
	UpdateProfilePictureToFalse(userId string)  error
}
