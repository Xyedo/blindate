package service

import (
	"time"

	"github.com/xyedo/blindate/pkg/domain/online"
	onlineEntities "github.com/xyedo/blindate/pkg/domain/online/entities"
)

func NewOnline(onlineRepo online.Repository) *Online {
	return &Online{
		onlineRepository: onlineRepo,
	}
}

type Online struct {
	onlineRepository online.Repository
}

func (o *Online) CreateNewOnline(userId string) error {
	err := o.onlineRepository.InsertNewOnline(onlineEntities.DTO{UserId: userId, LastOnline: time.Now(), IsOnline: false})
	if err != nil {
		return err
	}
	return nil
}
func (o *Online) GetOnline(userId string) (onlineEntities.DTO, error) {
	userOnline, err := o.onlineRepository.SelectOnline(userId)
	if err != nil {
		return onlineEntities.DTO{}, err
	}
	return userOnline, nil

}
func (o *Online) PutOnline(userId string, online bool) error {
	err := o.onlineRepository.UpdateOnline(userId, online)
	if err != nil {
		return err
	}
	return nil
}
