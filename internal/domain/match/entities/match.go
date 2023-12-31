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
