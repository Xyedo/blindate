package entities

import "time"

type Conversation struct {
	MatchId   string
	ChatRows  int64
	DayPass   int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64
}

