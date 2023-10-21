package optional

import (
	"database/sql"
	"database/sql/driver"
)

type Bool struct {
	Option[bool]
}

func (b Bool) Value() (driver.Value, error) {
	v, ok := b.Get()
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (b *Bool) Scan(value interface{}) error {
	sqlBool := sql.NullBool{}
	err := sqlBool.Scan(value)
	if err != nil {
		return err
	}

	if sqlBool.Valid {
		b.Set(sqlBool.Bool)
	}
	return nil
}
