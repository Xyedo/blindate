package entities

import (
	"strconv"

	geo "github.com/paulmach/go.geo"
	userentities "github.com/xyedo/blindate/internal/domain/user/entities"
)

func (matchUsers MatchUsers) CalculateDistance(user userentities.UserDetail) error {
	fUserLat, err := strconv.ParseFloat(user.Geog.Lat, 64)
	if err != nil {
		return err
	}

	fUserLng, err := strconv.ParseFloat(user.Geog.Lng, 64)
	if err != nil {
		return err
	}

	for i := range matchUsers {
		matchUserLat, err := strconv.ParseFloat(matchUsers[i].Geog.Lat, 64)
		if err != nil {
			return err
		}

		matchUserLng, err := strconv.ParseFloat(matchUsers[i].Geog.Lng, 64)
		if err != nil {
			return err
		}

		matchUsers[i].Distance = geo.
			NewPointFromLatLng(matchUserLat, matchUserLng).
			DistanceFrom(
				geo.NewPoint(fUserLat, fUserLng),
			)
	}

	return nil
}
