package pagination

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"
)

var ErrInvalidCursorFormat = errors.New("pagination: invalid cursor format")

type Pagination struct {
	Page  int
	Limit int
}

func (p Pagination) Offset() int {
	return p.Page*p.Limit - p.Limit
}

type Cursor struct {
	Id   string
	Date time.Time
}

func NewCursorFromBase64(b64 string) (Cursor, error) {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return Cursor{}, err
	}

	splitted := strings.Split(string(b), "-")
	if len(splitted) != 2 {
		return Cursor{}, ErrInvalidCursorFormat
	}
	parsedDate, err := time.Parse(time.RFC3339, splitted[1])
	if err != nil {
		return Cursor{}, ErrInvalidCursorFormat
	}

	return Cursor{
		Id:   splitted[0],
		Date: parsedDate,
	}, nil

}
