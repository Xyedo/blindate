package entities

import (
	geo "github.com/paulmach/go.geo"
	userentities "github.com/xyedo/blindate/internal/domain/user/entities"
)

func NewMatchUsers(requestUser userentities.UserDetail, matchUserDetails userentities.UserDetails, matchUserIdToMatchId map[string]string) []MatchUser {
	matchUsers := make([]MatchUser, 0, len(matchUserDetails))

	for _, matchUserDetail := range matchUserDetails {
		matchId, ok := matchUserIdToMatchId[matchUserDetail.UserId]
		if !ok {
			panic("should not happened")
		}

		matchUsers = append(matchUsers, MatchUser{
			MatchId:    matchId,
			UserDetail: matchUserDetail,
			Distance: geo.
				NewPointFromLatLng(matchUserDetail.Geog.Lat, matchUserDetail.Geog.Lng).
				DistanceFrom(
					geo.NewPoint(requestUser.Geog.Lat, requestUser.Geog.Lng),
				),
		})
	}

	return matchUsers
}
