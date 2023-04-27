package usecase

import (
	"errors"

	"github.com/xyedo/blindate/internal/security"
	"github.com/xyedo/blindate/pkg/authentication"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	"github.com/xyedo/blindate/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

func NewAuth(authRepo authentication.Repository, userRepo user.Repository, token *security.Jwt) *Auth {
	return &Auth{
		authRepo: authRepo,
		userRepo: userRepo,
		tokenSvc: token,
	}
}

type Auth struct {
	authRepo authentication.Repository
	userRepo user.Repository
	tokenSvc *security.Jwt
}

func (a *Auth) Login(email, password string) (accessToken string, refreshToken string, err error) {
	user, err := a.userRepo.GetUserByEmail(email)
	if err != nil {
		return
	}

	err = security.ComparePassword(user.Password, password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", "", apperror.Unauthorized(apperror.Payload{Error: err, Message: "password or email is invalid"})
		}
		return
	}
	
	accessToken = a.tokenSvc.GenerateAccessToken(user.ID)
	refreshToken = a.tokenSvc.GenerateRefreshToken(user.ID)

	err = a.authRepo.AddRefreshToken(refreshToken)
	if err != nil {
		return
	}

	return accessToken, refreshToken, err
}

func (a *Auth) RevalidateRefreshToken(refreshToken string) (string, error) {
	err := a.authRepo.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	id, err := a.tokenSvc.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	accessToken := a.tokenSvc.GenerateAccessToken(id)
	return accessToken, nil

}

func (a *Auth) Logout(refreshToken string) error {
	err := a.authRepo.VerifyRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	err = a.authRepo.DeleteRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	return nil
}
