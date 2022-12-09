package service

import (
	"time"

	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

func NewOnline(onlineRepo repository.Online) *Online {
	return &Online{
		onlineRepository: onlineRepo,
	}
}

type Online struct {
	onlineRepository repository.Online
}

func (o *Online) CreateNewOnline(userId string) error {
	err := o.onlineRepository.InsertNewOnline(domain.Online{UserId: userId, LastOnline: time.Now(), IsOnline: false})
	if err != nil {
		return err
	}
	return nil
}
func (o *Online) GetOnline(userId string) (domain.Online, error) {
	userOnline, err := o.onlineRepository.SelectOnline(userId)
	if err != nil {
		return domain.Online{}, err
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
