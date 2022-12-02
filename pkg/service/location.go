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
			return domain.ErrTooLongAccessingDB
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return ErrRefUserIdField
			}
			if pqErr.Code == "23505" {
				return domain.ErrUniqueConstraint23505
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
			return domain.ErrTooLongAccessingDB
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
			return nil, domain.ErrTooLongAccessingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		return nil, err
	}
	latlng := strings.TrimPrefix(location.Geog, "Point(")
	latlng = strings.TrimSuffix(latlng, ")")
	res := strings.Fields(latlng)

	return &domain.Location{
		UserId: location.UserId,
		Lat:    res[0],
		Lng:    res[1],
	}, nil
}

func latLngToGeog(lat, lng string) string {
	return fmt.Sprintf("Point(%s %s)", lat, lng)
}
