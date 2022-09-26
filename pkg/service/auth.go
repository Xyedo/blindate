package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
)

func NewAuth(auth domain.AuthRepository) *Auth {
	return &Auth{
		authRepo: auth,
	}
}

type Auth struct {
	authRepo domain.AuthRepository
}

func (a *Auth) AddRefreshToken(token string) error {
	_, err := a.authRepo.AddRefreshToken(token)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return domain.ErrDuplicateToken
			}
			log.Panicln(err.Error())

		}
		log.Panicln(err.Error())
		return err
	}
	return nil
}

func (a *Auth) VerifyRefreshToken(token string) error {
	err := a.authRepo.VerifyRefreshToken(token)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotMatchCredential
		}
		log.Panicln(err.Error())
	}
	return nil
}

func (a *Auth) DeleteRefreshToken(token string) error {
	rows, err := a.authRepo.DeleteRefreshToken(token)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		log.Panicln(err)
	}
	if rows == 0 {
		return domain.ErrNotMatchCredential
	}
	return nil
}
