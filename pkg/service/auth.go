package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

func NewAuth(authR repository.AuthRepo) *auth {
	return &auth{
		authRepo: authR,
	}
}

type auth struct {
	authRepo repository.AuthRepo
}

func (a *auth) AddRefreshToken(token string) error {
	_, err := a.authRepo.AddRefreshToken(token)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return domain.ErrUniqueConstraint23505
			}
		}
		return err
	}
	return nil
}

func (a *auth) VerifyRefreshToken(token string) error {
	err := a.authRepo.VerifyRefreshToken(token)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotMatchCredential
		}
		return err
	}
	return nil
}

func (a *auth) DeleteRefreshToken(token string) error {
	rows, err := a.authRepo.DeleteRefreshToken(token)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	if rows == 0 {
		return domain.ErrNotMatchCredential
	}
	return nil
}
