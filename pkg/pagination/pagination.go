package pagination

import (
	"encoding/base64"
	"errors"
	"fmt"
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

func (p Pagination) Next(hasNexts ...bool) *string {
	if len(hasNexts) > 0 && hasNexts[0] {
		return nil
	}
	r := fmt.Sprintf("?page=%d&limit=%d", p.Page+1, p.Limit)
	return &r
}

func (p Pagination) Prev(hasPrevs ...bool) *string {
	if len(hasPrevs) > 0 && hasPrevs[0] {
		return nil
	}

	if p.Page <= 1 {
		return nil
	}
	r := fmt.Sprintf("?page=%d&limit=%d", p.Page-1, p.Limit)
	return &r

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

func NewBase64FromCursor(cursor Cursor) string {
	s1 := cursor.Id
	s2 := cursor.Date.Format(time.RFC3339)

	s := strings.Join([]string{s1, s2}, "-")

	return base64.StdEncoding.EncodeToString([]byte(s))
}
