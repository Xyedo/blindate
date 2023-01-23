package location

import (
	locationEntity "github.com/xyedo/blindate/pkg/domain/location/entities"
	matchEntity "github.com/xyedo/blindate/pkg/domain/match/entities"
)

type Repository interface {
	InsertNewLocation(location *locationEntity.DAO) error
	UpdateLocation(location *locationEntity.DAO) error
	GetLocationByUserId(id string) (locationEntity.DAO, error)
	GetClosestUser(userId, geom string, limit int) ([]matchEntity.UserDTO, error)
}
