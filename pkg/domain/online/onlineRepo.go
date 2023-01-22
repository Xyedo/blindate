package online

import onlineEntities "github.com/xyedo/blindate/pkg/domain/online/entities"

type Repository interface {
	InsertNewOnline(on onlineEntities.DTO) error
	UpdateOnline(userId string, online bool) error
	SelectOnline(userId string) (onlineEntities.DTO, error)
}
