package online

import (
	onlineEntities "github.com/xyedo/blindate/pkg/domain/online/entities"
)

type Repository interface {
	Insert(onlineEntities.Online) error
	Update(string, bool) error
	SelectByUserId(string) (onlineEntities.Online, error)
}
