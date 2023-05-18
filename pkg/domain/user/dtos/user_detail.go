package userDTOs

import (
	"time"

	userEntities "github.com/xyedo/blindate/pkg/domain/user/entities"
)

type UserProfilePic struct {
	Id          string
	UserId      string
	Selected    bool
	PictureLink string
}

type User struct {
	ID         string
	FullName   string
	Alias      string
	ProfilePic []UserProfilePic
	Email      string
	Dob        time.Time
	Active     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func FromEntities(userEntity userEntities.User) User {
	user := User{
		ID:         userEntity.ID,
		FullName:   userEntity.FullName,
		Alias:      userEntity.Alias,
		Email:      userEntity.Email,
		ProfilePic: make([]UserProfilePic, 0, len(userEntity.ProfilePic)),
		Dob:        userEntity.Dob,
		Active:     userEntity.Active,
		CreatedAt:  userEntity.CreatedAt,
		UpdatedAt:  userEntity.UpdatedAt,
	}

	for _, pic := range userEntity.ProfilePic {
		user.ProfilePic = append(user.ProfilePic, UserProfilePic(pic))
	}

	return user
}
