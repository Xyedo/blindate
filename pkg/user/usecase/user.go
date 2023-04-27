package usecase

import (
	"path/filepath"
	"time"

	"github.com/xyedo/blindate/internal/security"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	"github.com/xyedo/blindate/pkg/user"
	userDTOs "github.com/xyedo/blindate/pkg/user/dtos"
	"github.com/xyedo/blindate/pkg/user/entities"
)

func NewUserUsecase(userRepo user.Repository) user.Usecase {
	return &userUC{
		userRepo: userRepo,
	}
}

type userUC struct {
	userRepo user.Repository
}

func (u *userUC) CreateUser(newUser userDTOs.RegisterUser) (string, error) {
	hashedPass, err := security.HashAndSalt(newUser.Password)
	if err != nil {
		return "", err
	}
	newUser.Password = hashedPass

	userId, err := u.userRepo.InsertUser(newUser)
	if err != nil {
		return "", err
	}

	return userId, nil
}
func (u *userUC) GetUserById(userDetail userDTOs.GetUserDetail) (entities.User, error) {
	user, err := u.userRepo.GetUserById(userDetail.Id)
	if err != nil {
		return entities.User{}, err
	}

	if userDetail.ProfilePicture {
		profilePicture, err := u.userRepo.SelectProfilePicture(userDetail.Id, userDTOs.ProfilePictureQuery{})
		if err != nil {
			return entities.User{}, err
		}
		user.ProfilePic = profilePicture
	}

	return user, nil
}

func (u *userUC) UpdateUser(updateUser userDTOs.UpdateUser) error {
	olduser, err := u.userRepo.GetUserById(updateUser.Id)
	if err != nil {
		return err
	}
	if updateUser.NewPassword.Present() && updateUser.OldPassword.Present() {

		err := security.ComparePassword(olduser.Password, updateUser.OldPassword.MustGet())
		if err != nil {
			return apperror.Unauthorized(apperror.Payload{Error: err, Message: "invalid email or password"})
		}

		hashedPass, err := security.HashAndSalt(updateUser.NewPassword.MustGet())
		if err != nil {
			return err
		}
		olduser.Password = hashedPass
	}

	updateUser.FullName.If(func(userFullName string) {
		olduser.FullName = userFullName
	})

	updateUser.Alias.If(func(userAlias string) {
		olduser.Alias = userAlias
	})

	updateUser.Email.If(func(userEmail string) {
		olduser.Active = false
		olduser.Email = userEmail
	})

	updateUser.Dob.If(func(userDOB time.Time) {
		olduser.Dob = userDOB
	})

	err = u.userRepo.UpdateUser(olduser)
	if err != nil {
		return err
	}

	// if updateUser.Alias.Present() || updateUser.FullName.Present() {
	// 	event.ProfileUpdated.Trigger(event.ProfileUpdatedPayload{
	// 		UserId: userId,
	// 	})
	// }

	return nil
}

func (u *userUC) CreateNewProfilePic(profilePicture userDTOs.RegisterProfilePicture) (string, error) {
	profilePictures, err := u.userRepo.SelectProfilePicture(profilePicture.UserId, userDTOs.ProfilePictureQuery{})
	if err != nil {
		return "", err
	}
	if len(profilePictures) >= 5 {
		return "", apperror.UnprocessableEntity(apperror.PayloadMap{ErrorMap: map[string][]string{"profilePicture": {"maximal profile pics is 5"}}})
	}
	if profilePicture.Selected {
		err := u.userRepo.UpdateProfilePictureToFalse(profilePicture.UserId)
		if err != nil {
			return "", err
		}
	}

	profilePicture.PictureLink = filepath.Base(profilePicture.PictureLink)
	id, err := u.userRepo.CreateProfilePicture(profilePicture.UserId, profilePicture.PictureLink, profilePicture.Selected)
	if err != nil {
		return "", err
	}

	// if profilePicture.Selected {
	// 	event.ProfileUpdated.Trigger(event.ProfileUpdatedPayload{
	// 		UserId: profPicParam.UserId,
	// 	})
	// }
	return id, nil
}
