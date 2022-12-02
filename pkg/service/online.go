package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

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
			if pqErr.Code == "23503" {
				return ErrRefUserIdField
			}
			if pqErr.Code == "23505" {
				return domain.ErrUniqueConstraint23505
			}
			return err
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccessingDB
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
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccessingDB
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
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccessingDB
		}
		return err
	}
	return nil
}
