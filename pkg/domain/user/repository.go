package user

import (
	userDTOs "github.com/xyedo/blindate/pkg/domain/user/dtos"
	userEntities "github.com/xyedo/blindate/pkg/domain/user/entities"
)

type Repository interface {
	InsertUser(user userDTOs.RegisterUser) (string, error)
	GetUserById(id string) (userEntities.User, error)
	GetUserByEmail(email string) (userEntities.User, error)
	UpdateUser(user userEntities.User) error
	CreateProfilePicture(userId, pictureRef string, selected bool) (string, error)
	SelectProfilePicture(userId string, params userDTOs.ProfilePictureQuery) ([]userEntities.UserProfilePic, error)
	UpdateProfilePictureToFalse(userId string) error
}
