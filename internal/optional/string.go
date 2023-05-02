package optional

import (
	"database/sql"
	"database/sql/driver"
)

type String struct {
	Option[string]
}

func (s String) Value() (driver.Value, error) {
	str, ok := s.Get()
	if !ok {
		return nil, nil
	}
	return str, nil
}

func (s *String) Scan(value interface{}) error {
	sqlStr := sql.NullString{}
	err := sqlStr.Scan(value)
	if err != nil {
		return err
	}

	if sqlStr.Valid {
		s.Set(sqlStr.String)
	}
	return nil
}
