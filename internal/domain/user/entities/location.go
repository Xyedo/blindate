package entities

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Geography struct {
	Lat float64
	Lng float64
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

	var err error

	g.Lat, err = strconv.ParseFloat(geogFields[0], 64)
	if err != nil {
		return err
	}

	g.Lng, err = strconv.ParseFloat(geogFields[1], 64)
	if err != nil {
		return err
	}

	return nil
}

func (g Geography) Value() (driver.Value, error) {
	return fmt.Sprintf("POINT(%f %f)", g.Lat, g.Lng), nil
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
