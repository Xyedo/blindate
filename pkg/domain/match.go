package domain

import (
	"errors"
	"strconv"
	"time"
)

type BigUser struct {
	UserId           string       ` json:"userId"`
	Alias            string       ` json:"alias"`
	Dob              time.Time    ` json:"dob"`
	Gender           *string      `json:"gender"`
	FromLoc          *string      `json:"fromLoc"`
	Height           *int         `json:"height"`
	EducationLevel   *string      `json:"educationLevel"`
	Drinking         *string      `json:"drinking"`
	Smoking          *string      `json:"smoking"`
	RelationshipPref *string      `json:"relationshipPref"`
	LookingFor       *string      `json:"lookingFor"`
	Zodiac           *string      `json:"zodiac"`
	Kids             *int         `json:"kids"`
	Work             *string      `json:"work"`
	BioId            *string      `json:"bioId"`
	Bio              *string      `json:"bio" db:"bio"`
	Hobbies          []Hobbie     `json:"hobbies"`
	MovieSeries      []MovieSerie `json:"movieSeries"`
	Travels          []Travel     `json:"travels"`
	Sports           []Sport      `json:"sports"`
}
type MatchUser struct {
	MatchId string `json:"matchId"`
	BigUser
}

type MatchStatus string

const (
	Unknown   MatchStatus = "unknown"
	Requested MatchStatus = "requested"
	Declined  MatchStatus = "declined"
	Accepted  MatchStatus = "accepted"
)

var ErrInvalidMatchStatusFormat = errors.New("invalid MatchStatus format")

func (m *MatchStatus) UnmarshalJSON(jsonValue []byte) error {
	unquotedJsonV, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidMatchStatusFormat
	}
	*m = MatchStatus(unquotedJsonV)
	return nil
}
func (m *MatchStatus) MarshalJSON() ([]byte, error) {
	if *m == Unknown {
		*m = ""
	}
	jsonValue := string(*m)
	quotedJSONvalue := strconv.Quote(jsonValue)
	return []byte(quotedJSONvalue), nil
}
