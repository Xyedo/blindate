package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

func NewAuth(authR repository.Auth) auth {
	return auth{
		authRepo: authR,
	}
}

type auth struct {
	authRepo repository.Auth
}

func (a auth) AddRefreshToken(token string) error {
	_, err := a.authRepo.AddRefreshToken(token)
	if err != nil {
		return a.wrapError(err)
	}
	return nil
}

func (a auth) VerifyRefreshToken(token string) error {
	err := a.authRepo.VerifyRefreshToken(token)
	if err != nil {
		return a.wrapError(err)
	}
	return nil
}

func (a auth) DeleteRefreshToken(token string) error {
	rows, err := a.authRepo.DeleteRefreshToken(token)
	if err != nil {
		return a.wrapError(err)
	}
	if rows == 0 {
		return domain.WrapError(err, domain.ErrNotMatchCredential)
	}
	return nil
}

func (auth) wrapError(err error) error {
	var pqErr *pq.Error
	switch {
	case errors.Is(err, context.Canceled):
		return domain.WrapError(err, domain.ErrTooLongAccessingDB)
	case errors.Is(err, sql.ErrNoRows):
		return domain.WrapError(err, domain.ErrNotMatchCredential)
	case errors.As(err, &pqErr) && pqErr.Code == "23505":
		return domain.WrapErrorWithMsg(err, domain.ErrUniqueConstraint23505, "token is already taken, please try again")
	default:
		return err
	}
}
