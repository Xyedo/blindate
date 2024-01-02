package entities

import "github.com/xyedo/blindate/pkg/pagination"

type FindUserMatchByStatus struct {
	UserId     string
	Statuses   []MatchStatus
	Pagination pagination.Pagination
}
type GetMatchOption struct {
	PessimisticLocking bool
}

func (matchs Matchs) ToUserIds(requestId string) ([]string, map[string]string) {
	userIdToMatchId := make(map[string]string, 0)
	userIds := make([]string, 0, len(matchs))

	for _, match := range matchs {
		if requestId == match.RequestFrom {
			userIds = append(userIds, match.RequestTo)
			userIdToMatchId[match.RequestTo] = match.Id
		}

		if requestId == match.RequestTo {
			userIds = append(userIds, match.RequestFrom)
			userIdToMatchId[match.RequestFrom] = match.Id
		}
	}

	return userIds, userIdToMatchId
}
