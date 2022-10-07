package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

type Online interface {
	CreateNewOnline(userId string) error
	PutOnline(userId string, online bool) error
	GetOnline(userId string) (*domain.Online, error)
}

func NewOnline(onlineRepo repository.Online) *online {
	return &online{
		onlineRepository: onlineRepo,
	}
}

type online struct {
	onlineRepository repository.Online
}

func (o *online) CreateNewOnline(userId string) error {
	err := o.onlineRepository.InsertNewOnline(&domain.Online{UserId: userId, LastOnline: time.Now(), IsOnline: false})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if strings.Contains(pqErr.Constraint, "user_id") {
				if pqErr.Code == "23503" {
					return domain.ErrResourceNotFound
				}
				if pqErr.Code == "23505" {
					return ErrUniqueConstrainUserId
				}
			}
			return err
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	return nil
}
func (o *online) GetOnline(userId string) (*domain.Online, error) {
	userOnline, err := o.onlineRepository.SelectOnline(userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		return nil, err
	}
	return userOnline, nil

}
func (o *online) PutOnline(userId string, online bool) error {
	err := o.onlineRepository.UpdateOnline(userId, online)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrResourceNotFound
		}
		return err
	}
	return nil
}
