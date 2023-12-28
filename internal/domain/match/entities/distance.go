package entities

import (
	geo "github.com/paulmach/go.geo"
	userentities "github.com/xyedo/blindate/internal/domain/user/entities"
)

func (matchUsers MatchUsers) CalculateDistance(user userentities.UserDetail) error {
	for i := range matchUsers {
		matchUsers[i].Distance = geo.
			NewPointFromLatLng(matchUsers[i].Geog.Lat, matchUsers[i].Geog.Lng).
			DistanceFrom(
				geo.NewPoint(user.Geog.Lat, user.Geog.Lng),
			)
	}

	return nil
}
