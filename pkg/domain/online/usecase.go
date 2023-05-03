package online

import (
	"github.com/xyedo/blindate/pkg/common/event"
	onlineEntities "github.com/xyedo/blindate/pkg/domain/online/entities"
)

type Usecase interface {
	Create(string) error
	GetByUserId(string, string) (onlineEntities.Online, error)
	HandleUserChangeConnection(payload event.ConnectionPayload)
}
