package optional

import (
	"database/sql"
	"database/sql/driver"
)

func NewInt64(values ...int64) Int64 {
	var b Int64
	b.set = true

	if len(values) > 0 {
		b.value = &values[0]
	}

	return b
}

type Int64 struct {
	Option[int64]
}

func (i Int64) Value() (driver.Value, error) {
	v, ok := i.Get()
	if !i.IsSet() || !ok {
		return nil, nil
	}
	return v, nil
}

func (b *Int64) Scan(value interface{}) error {
	sqlInt64 := sql.NullInt64{}
	err := sqlInt64.Scan(value)
	if err != nil {
		return err
	}

	if sqlInt64.Valid {
		b.Set(sqlInt64.Int64)
	}
	return nil
}
func NewInt32(values ...int32) Int32 {
	var b Int32
	b.set = true

	if len(values) > 0 {
		b.value = &values[0]
	}

	return b
}

type Int32 struct {
	Option[int32]
}

func (i Int32) Value() (driver.Value, error) {
	v, ok := i.Get()
	if !i.IsSet() || !ok {
		return nil, nil
	}
	return int64(v), nil
}

func (i *Int32) Scan(value interface{}) error {
	sqlInt32 := sql.NullInt32{}
	err := sqlInt32.Scan(value)
	if err != nil {
		return err
	}

	if sqlInt32.Valid {
		i.Set(sqlInt32.Int32)
	}
	return nil
}
func NewInt16(values ...int16) Int16 {
	var b Int16
	b.set = true

	if len(values) > 0 {
		b.value = &values[0]
	}

	return b
}

type Int16 struct {
	Option[int16]
}

func (i Int16) Value() (driver.Value, error) {
	v, ok := i.Get()
	if !ok {
		return nil, nil
	}
	return int64(v), nil
}

func (i *Int16) Scan(value interface{}) error {
	sqlInt16 := sql.NullInt16{}
	err := sqlInt16.Scan(value)
	if err != nil {
		return err
	}

	if sqlInt16.Valid {
		i.Set(sqlInt16.Int16)
	}
	return nil
}
