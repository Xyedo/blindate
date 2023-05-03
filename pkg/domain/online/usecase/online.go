package usecase

import (
	"time"

	"github.com/xyedo/blindate/pkg/common/event"
	"github.com/xyedo/blindate/pkg/domain/online"
	onlineEntities "github.com/xyedo/blindate/pkg/domain/online/entities"
)

func New(onlineRepo online.Repository) online.Usecase {
	return &onlineUC{
		onlineRepo: onlineRepo,
	}
}

type onlineUC struct {
	onlineRepo online.Repository
}

// Create implements online.Usecase
func (uc *onlineUC) Create(userId string) error {
	err := uc.onlineRepo.Insert(onlineEntities.Online{
		UserId:     userId,
		LastOnline: time.Now(),
		IsOnline:   true,
	})
	if err != nil {
		return err
	}

	return nil
}

// GetByUserId implements online.Usecase
func (uc *onlineUC) GetByUserId(requestId, userId string) (onlineEntities.Online, error) {
	foundOnline, err := uc.onlineRepo.SelectByUserId(userId)
	if err != nil {
		return onlineEntities.Online{}, err
	}

	return foundOnline, nil
}

func (uc *onlineUC) HandleUserChangeConnection(payload event.ConnectionPayload) {
	uc.update(payload.UserId, payload.Online)
}

// update implements online.Usecase
func (uc *onlineUC) update(userId string, online bool) error {
	err := uc.onlineRepo.Update(userId, online)
	if err != nil {
		return err
	}

	return nil
}
