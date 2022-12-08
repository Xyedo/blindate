package service

import (
	"errors"

	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

func NewAuth(authR repository.Auth, userR repository.User, tokenSvc *Jwt) *Auth {
	return &Auth{
		authRepo: authR,
		userRepo: userR,
		tokenSvc: tokenSvc,
	}
}

type Auth struct {
	authRepo repository.Auth
	userRepo repository.User
	tokenSvc *Jwt
}

func (a *Auth) Login(email, password string) (accessToken string, refreshToken string, err error) {
	user, err := a.userRepo.GetUserByEmail(email)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", "", common.WrapError(err, common.ErrNotMatchCredential)
		}
		return
	}
	accessToken, err = a.tokenSvc.GenerateAccessToken(user.ID)
	if err != nil {
		panic(err)
	}
	refreshToken, err = a.tokenSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		panic(err)
	}
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
	accessToken, err := a.tokenSvc.GenerateAccessToken(id)
	if err != nil {
		panic(err)
	}
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
