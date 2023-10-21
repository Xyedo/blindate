package optional

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type Time struct {
	Option[time.Time]
}

func (t Time) Value() (driver.Value, error) {
	v, ok := t.Get()
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (t *Time) Scan(value interface{}) error {
	sqlTime := sql.NullTime{}
	err := sqlTime.Scan(value)
	if err != nil {
		return err
	}

	if sqlTime.Valid {
		t.Set(sqlTime.Time)
	}

	return nil
}
