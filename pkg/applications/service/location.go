package service

import (
	"fmt"
	"strings"

	"github.com/xyedo/blindate/pkg/domain/location"
	locationEntity "github.com/xyedo/blindate/pkg/domain/location/entities"
)

func NewLocation(locationRepo location.Repository) *Location {
	return &Location{
		locationRepo: locationRepo,
	}
}

type Location struct {
	locationRepo location.Repository
}

func (l *Location) CreateNewLocation(location *locationEntity.DTO) error {
	locationDAO := &locationEntity.DAO{
		UserId: location.UserId,
		Geog:   latLngToGeog(location.Lat, location.Lng),
	}
	err := l.locationRepo.InsertNewLocation(locationDAO)
	if err != nil {
		return err
	}

	return nil

}

func (l *Location) UpdateLocation(userId string, changeLat, changeLng *string) error {
	location, err := l.GetLocation(userId)
	if err != nil {
		return err
	}
	if changeLat != nil {
		location.Lat = *changeLat
	}
	if changeLng != nil {
		location.Lng = *changeLng
	}
	locationDAO := &locationEntity.DAO{
		UserId: location.UserId,
		Geog:   latLngToGeog(location.Lat, location.Lng),
	}
	err = l.locationRepo.UpdateLocation(locationDAO)
	if err != nil {
		return err
	}
	return nil
}

func (l *Location) GetLocation(id string) (locationEntity.DTO, error) {
	location, err := l.locationRepo.GetLocationByUserId(id)
	if err != nil {
		return locationEntity.DTO{}, err
	}
	latlng := strings.TrimPrefix(location.Geog, "POINT(")
	latlng = strings.TrimSuffix(latlng, ")")
	res := strings.Fields(latlng)

	return locationEntity.DTO{
		UserId: location.UserId,
		Lat:    res[0],
		Lng:    res[1],
	}, nil
}

func latLngToGeog(lat, lng string) string {
	return fmt.Sprintf("POINT(%s %s)", lat, lng)
}
