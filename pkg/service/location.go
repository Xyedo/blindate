package service

import (
	"fmt"
	"strings"

	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/domain/entity"
	"github.com/xyedo/blindate/pkg/repository"
)

func NewLocation(locationRepo repository.Location) *Location {
	return &Location{
		locationRepo: locationRepo,
	}
}

type Location struct {
	locationRepo repository.Location
}

func (l *Location) CreateNewLocation(location *domain.Location) error {
	locationEntity := &entity.Location{
		UserId: location.UserId,
		Geog:   latLngToGeog(location.Lat, location.Lng),
	}
	err := l.locationRepo.InsertNewLocation(locationEntity)
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
	locationEntity := &entity.Location{
		UserId: location.UserId,
		Geog:   latLngToGeog(location.Lat, location.Lng),
	}
	err = l.locationRepo.UpdateLocation(locationEntity)
	if err != nil {
		return err
	}
	return nil
}

func (l *Location) GetLocation(id string) (domain.Location, error) {
	location, err := l.locationRepo.GetLocationByUserId(id)
	if err != nil {
		return domain.Location{}, err
	}
	latlng := strings.TrimPrefix(location.Geog, "Point(")
	latlng = strings.TrimSuffix(latlng, ")")
	res := strings.Fields(latlng)

	return domain.Location{
		UserId: location.UserId,
		Lat:    res[0],
		Lng:    res[1],
	}, nil
}

func latLngToGeog(lat, lng string) string {
	return fmt.Sprintf("Point(%s %s)", lat, lng)
}
