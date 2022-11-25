package domain

import (
	"errors"
	"strconv"
)

type MatchStatus string

const (
	Unknown   MatchStatus = "unknown"
	Requested MatchStatus = "requested"
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
