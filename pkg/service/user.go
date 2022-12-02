package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

// var ErrMaxProfilePicture =
var ErrMaxProfilePicture = domain.WrapWithNewError(errors.New("excedeed profile picture constraint"), http.StatusUnprocessableEntity, "maximal profile pics is 5")

func NewUser(userRepo repository.User) user {
	return user{
		userRepository: userRepo,
	}
}

type user struct {
	userRepository repository.User
}

func (u user) CreateUser(newUser *domain.User) error {
	hashedPass, err := hashAndSalt(newUser.Password)
	if err != nil {
		return err
	}
	newUser.HashedPassword = hashedPass
	newUser.Password = ""

	err = u.userRepository.InsertUser(newUser)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return domain.WrapErrorWithMsg(err, domain.ErrUniqueConstraint23505, "email already taken")
			}
			return pqErr
		}
		if errors.Is(err, context.Canceled) {
			return domain.WrapError(err, domain.ErrTooLongAccessingDB)
		}
		return err
	}

	return nil
}
func (u user) GetUserById(id string) (*domain.User, error) {
	user, err := u.userRepository.GetUserById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.WrapError(err, domain.ErrResourceNotFound)
		}
		if errors.Is(err, context.Canceled) {
			return nil, domain.WrapError(err, domain.ErrTooLongAccessingDB)
		}
		return nil, err
	}

	return user, nil
}
func (u user) GetUserIdWithProfilePics(id string) (*domain.User, error) {
	user, err := u.GetUserById(id)
	if err != nil {
		return nil, err
	}

	profPics, err := u.userRepository.SelectProfilePicture(id, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.WrapError(err, domain.ErrResourceNotFound)
		}
		if errors.Is(err, context.Canceled) {
			return nil, domain.WrapError(err, domain.ErrTooLongAccessingDB)
		}
		return nil, err
	}
	user.ProfilePic = profPics
	return user, nil
}
func (u user) VerifyCredential(email, password string) (string, error) {
	user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", domain.WrapError(err, domain.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.WrapError(err, domain.ErrNotMatchCredential)
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", domain.WrapError(err, domain.ErrNotMatchCredential)
		}
		return "", err
	}
	return user.ID, nil

}
func (u user) UpdateUser(user *domain.User) error {
	if user.Password != "" {
		hashedPass, err := hashAndSalt(user.Password)
		if err != nil {
			return err
		}
		user.HashedPassword = hashedPass
		user.Password = ""
	}
	err := u.userRepository.UpdateUser(user)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.WrapError(err, domain.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.WrapError(err, domain.ErrResourceNotFound)
		}
		return err
	}
	return nil
}

func (u user) CreateNewProfilePic(profPicParam domain.ProfilePicture) (string, error) {
	profPics, err := u.userRepository.SelectProfilePicture(profPicParam.UserId, nil)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", domain.WrapError(err, domain.ErrTooLongAccessingDB)
		}
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.WrapErrorWithMsg(err, domain.ErrRefNotFound23503, "profile picture not found")
		}
		return "", err
	}
	if len(profPics) >= 5 {
		return "", ErrMaxProfilePicture
	}
	if profPicParam.Selected {
		_, err := u.userRepository.ProfilePicSelectedToFalse(profPicParam.UserId)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return "", domain.WrapError(err, domain.ErrTooLongAccessingDB)
			}
			return "", err
		}
	}
	profPicParam.PictureLink = filepath.Base(profPicParam.PictureLink)
	id, err := u.userRepository.CreateProfilePicture(profPicParam.UserId, profPicParam.PictureLink, profPicParam.Selected)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", domain.WrapError(err, domain.ErrTooLongAccessingDB)
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return "", domain.WrapErrorWithMsg(err, domain.ErrRefNotFound23503, "profile picture not found")
			}
			return "", pqErr
		}
		return "", err
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
