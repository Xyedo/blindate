package location

import (
	locationEntities "github.com/xyedo/blindate/pkg/domain/location/entities"
)

type Repository interface {
	InsertLocation(locationEntities.Location) error
	PatchLocation(userId, location string) error
	GetLocationByUserId(string) (locationEntities.Location, error)
	
}
