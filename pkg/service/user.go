package service

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain/user"
	userEntity "github.com/xyedo/blindate/pkg/domain/user/entities"
	"github.com/xyedo/blindate/pkg/event"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(userRepo user.Repository) *User {
	return &User{
		userRepository: userRepo,
	}
}

type User struct {
	userRepository user.Repository
}

func (u *User) CreateUser(newUser userEntity.Register) (string, error) {
	hashedPass, err := hashAndSalt(newUser.Password)
	if err != nil {
		return "", err
	}
	newUser.Password = hashedPass

	userId, err := u.userRepository.InsertUser(newUser)
	if err != nil {
		return "", err
	}

	return userId, nil
}
func (u *User) GetUserById(id string) (userEntity.FullDTO, error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		return userEntity.FullDTO{}, err
	}
	return user, nil
}
func (u *User) GetUserByIdWithSelectedProfPic(id string) (userEntity.FullDTO, error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		return userEntity.FullDTO{}, err
	}

	profPics, err := u.userRepository.SelectProfilePicture(id, nil)
	if err != nil {
		return userEntity.FullDTO{}, err
	}
	user.ProfilePic = profPics
	return user, nil
}

func (u *User) UpdateUser(userId string, updateUser userEntity.Update) error {
	olduser, err := u.userRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	if updateUser.NewPassword != nil && updateUser.OldPassword != nil {
		err = bcrypt.CompareHashAndPassword([]byte(olduser.Password), []byte(*updateUser.OldPassword))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return common.WrapError(err, common.ErrNotMatchCredential)
			}
			return err
		}

		hashedPass, err := hashAndSalt(*updateUser.NewPassword)
		if err != nil {
			return err
		}
		olduser.Password = hashedPass
	}

	if updateUser.FullName != nil {
		olduser.FullName = *updateUser.FullName
	}

	if updateUser.Alias != nil {
		olduser.Alias = *updateUser.Alias
	}

	if updateUser.Email != nil {
		olduser.Active = false
		olduser.Email = *updateUser.Email
	}
	if updateUser.Dob != nil {
		olduser.Dob = *updateUser.Dob
	}
	err = u.userRepository.UpdateUser(olduser)
	if err != nil {
		return err
	}
	if updateUser.Alias != nil || updateUser.FullName != nil {
		event.ProfileUpdated.Trigger(event.ProfileUpdatedPayload{
			UserId: userId,
		})
	}

	return nil
}

func (u *User) CreateNewProfilePic(profPicParam userEntity.ProfilePic) (string, error) {
	profPics, err := u.userRepository.SelectProfilePicture(profPicParam.UserId, nil)
	if err != nil {
		return "", err
	}
	if len(profPics) >= 5 {
		return "", common.WrapWithNewError(common.ErrMaxProfilePicture, http.StatusUnprocessableEntity, "maximal profile pics is 5")
	}
	if profPicParam.Selected {
		_, err := u.userRepository.ProfilePicSelectedToFalse(profPicParam.UserId)
		if err != nil {
			return "", err
		}
	}
	profPicParam.PictureLink = filepath.Base(profPicParam.PictureLink)
	id, err := u.userRepository.CreateProfilePicture(profPicParam.UserId, profPicParam.PictureLink, profPicParam.Selected)
	if err != nil {
		return "", err
	}
	if profPicParam.Selected {
		event.ProfileUpdated.Trigger(event.ProfileUpdatedPayload{
			UserId: profPicParam.UserId,
		})
	}
	return id, nil
}
func hashAndSalt(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPass), err
}
