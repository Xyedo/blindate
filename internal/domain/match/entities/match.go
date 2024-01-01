package entities

type FindUserMatchByStatus struct {
	UserId   string
	Statuses []MatchStatus
	Limit    int
	Page     int
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
