package matchEntity

import (
	"errors"
	"strconv"
)

type Status string

const (
	Unknown   Status = "unknown"
	Requested Status = "requested"
	Declined  Status = "declined"
	Accepted  Status = "accepted"
)

var ErrInvalidMatchStatusFormat = errors.New("invalid MatchStatus format")

func (m *Status) UnmarshalJSON(jsonValue []byte) error {
	unquotedJsonV, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidMatchStatusFormat
	}
	*m = Status(unquotedJsonV)
	return nil
}
func (m *Status) MarshalJSON() ([]byte, error) {
	if *m == Unknown {
		*m = ""
	}
	jsonValue := string(*m)
	quotedJSONvalue := strconv.Quote(jsonValue)
	return []byte(quotedJSONvalue), nil
}
