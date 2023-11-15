package entities

type FindUserMatchByStatus struct {
	UserId string
	Status MatchStatus
	Limit  int
	Page   int
}
