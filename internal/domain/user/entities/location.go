package entities

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

type Geography struct {
	Lat string
	Lng string
}

func (g *Geography) Scan(v any) error {
	geogStr, ok := v.(string)
	if !ok {
		return errors.New("invalid data-type")
	}

	geogStr = strings.TrimPrefix(geogStr, "POINT(")
	geogStr = strings.TrimSuffix(geogStr, ")")

	geogFields := strings.Fields(geogStr)
	if len(geogFields) < 2 {
		return errors.New("fields is less than 2")
	}

	g.Lat = geogFields[0]
	g.Lng = geogFields[1]

	return nil
}

func (g Geography) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%s %s)", g.Lat, g.Lng), nil
}

type GetLocationOption struct {
	PessimisticLocking bool
}

type CreateLocation struct {
	UserId string
	Geog   Geography
}

type FindClosestUser struct {
	UserId string
	Geog   Geography
	Page   int
	Limit  int
}

type ClosestUser struct {
	UserId   string
	Distance string
}
