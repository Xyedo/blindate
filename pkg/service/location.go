package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
	"github.com/xyedo/blindate/pkg/repository"
)

type Location interface {
	CreateNewLocation(location *domain.Location) error
	UpdateLocation(location *domain.Location) error
	GetLocation(id string) (*domain.Location, error)
}

func NewLocation(locationRepo repository.Location) *location {
	return &location{
		locationRepo: locationRepo,
	}
}

type location struct {
	locationRepo repository.Location
}

func (l *location) CreateNewLocation(location *domain.Location) error {
	locationEntity := &entity.Location{
		UserId: location.UserId,
		Geog:   latLngToGeog(location.Lat, location.Lng),
	}
	rows, err := l.locationRepo.InsertNewLocation(locationEntity)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return ErrUserIdField
			}
		}
		return err
	}
	if rows == 0 {
		panic(rows)
	}
	return nil

}

func (l *location) UpdateLocation(location *domain.Location) error {
	locationEntity := &entity.Location{
		UserId: location.UserId,
		Geog:   latLngToGeog(location.Lat, location.Lng),
	}
	rows, err := l.locationRepo.UpdateLocation(locationEntity)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	if rows == 0 {
		return domain.ErrResourceNotFound
	}
	return nil
}

func (l *location) GetLocation(id string) (*domain.Location, error) {
	location, err := l.locationRepo.GetLocationByUserId(id)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		return nil, err
	}
	latlng := strings.TrimPrefix(location.Geog, "Point(")
	latlng = strings.TrimSuffix(latlng, ")")
	res := strings.Split(latlng, " ")

	return &domain.Location{
		UserId: location.UserId,
		Lat:    res[0],
		Lng:    res[1],
	}, nil
}

func latLngToGeog(lat, lng string) string {
	return fmt.Sprintf("Point(%s %s)", lat, lng)
}
