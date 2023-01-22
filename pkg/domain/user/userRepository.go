package user

import (
	userEntities "github.com/xyedo/blindate/pkg/domain/user/entities"
)

type ProfilePicQuery struct {
	Selected *bool
}

type Repository interface {
	InsertUser(user userEntities.Register) (string, error)
	GetUserById(id string) (userEntities.FullDTO, error)
	GetUserByEmail(email string) (userEntities.FullDTO, error)
	UpdateUser(user userEntities.FullDTO) error
	CreateProfilePicture(userId, pictureRef string, selected bool) (string, error)
	SelectProfilePicture(userId string, params *ProfilePicQuery) ([]userEntities.ProfilePic, error)
	ProfilePicSelectedToFalse(userId string) (int64, error)
}
