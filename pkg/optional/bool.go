package optional

import (
	"database/sql"
	"database/sql/driver"
)

func NewBool(values ...bool) Bool {
	var b Bool
	b.set = true

	if len(values) > 0 {
		b.value = &values[0]
	}

	return b
}

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
