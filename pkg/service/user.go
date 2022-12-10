package service

import (
	"errors"
	"net/http"
	"path/filepath"

	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/event"
	"github.com/xyedo/blindate/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

func NewUser(userRepo repository.User) *User {
	return &User{
		userRepository: userRepo,
	}
}

type User struct {
	userRepository repository.User
}

func (u *User) CreateUser(newUser *domain.User) error {
	hashedPass, err := hashAndSalt(newUser.Password)
	if err != nil {
		return err
	}
	newUser.HashedPassword = hashedPass
	newUser.Password = ""

	err = u.userRepository.InsertUser(newUser)
	if err != nil {
		return err
	}

	return nil
}
func (u *User) GetUserById(id string) (domain.User, error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
func (u *User) GetUserByIdWithSelectedProfPic(id string) (domain.User, error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		return domain.User{}, err
	}

	profPics, err := u.userRepository.SelectProfilePicture(id, nil)
	if err != nil {
		return domain.User{}, err
	}
	user.ProfilePic = profPics
	return user, nil
}

func (u *User) UpdateUser(userId string, updateUser domain.UpdateUser) error {
	olduser, err := u.userRepository.GetUserById(userId)
	if err != nil {
		return err
	}
	if updateUser.NewPassword != nil && updateUser.OldPassword != nil {
		err = bcrypt.CompareHashAndPassword([]byte(olduser.HashedPassword), []byte(*updateUser.OldPassword))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return common.WrapError(err, common.ErrNotMatchCredential)
			}
			return err
		}

		olduser.Password = *updateUser.NewPassword
		hashedPass, err := hashAndSalt(olduser.Password)
		if err != nil {
			return err
		}
		olduser.HashedPassword = hashedPass
		olduser.Password = ""
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

func (u *User) CreateNewProfilePic(profPicParam domain.ProfilePicture) (string, error) {
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
