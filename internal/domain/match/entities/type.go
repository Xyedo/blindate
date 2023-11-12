package entities

type MatchStatus string

const (
	MatchStatusUnknown   MatchStatus = "UNKNOWN"
	MatchStatusRequested MatchStatus = "REQUESTED"
	MatchStatusAccepted  MatchStatus = "ACCEPTED"
	MatchStatusDeclined  MatchStatus = "DECLINED"
)
