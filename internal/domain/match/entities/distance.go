package entities

import (
	geo "github.com/paulmach/go.geo"
	userentities "github.com/xyedo/blindate/internal/domain/user/entities"
)

func NewMatchUsers(matchGeo userentities.Geography, userDetails userentities.UserDetails) []MatchUser {
	matchUsers := make([]MatchUser, 0, len(userDetails))

	for _, userDetail := range userDetails {
		matchUsers = append(matchUsers, MatchUser{
			UserDetail: userDetail,
			Distance: geo.
				NewPointFromLatLng(userDetail.Geog.Lat, userDetail.Geog.Lng).
				DistanceFrom(
					geo.NewPoint(matchGeo.Lat, matchGeo.Lng),
				),
		})
	}

	return matchUsers
}
