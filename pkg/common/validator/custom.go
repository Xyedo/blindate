package validator

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

func ValidDob(value any) error {
	date, ok := value.(time.Time)
	if !ok {
		return errors.New("not valid data type")
	}
	curr := time.Now()

	youngest := time.Date(
		curr.Year()-18,
		curr.Month(),
		curr.Day(),
		curr.Hour(),
		curr.Minute(),
		curr.Second(),
		curr.Nanosecond(),
		curr.Location(),
	)
	oldest := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	if !(date.Before(youngest) && date.After(oldest)) {
		return errors.New("invalid day of birth")
	}
	return nil
}

func ValidUsername(value any) error {
	var s string
	v, ok := value.(interface{ Value() (driver.Value, error) })
	if !ok {
		v, err := v.Value()
		if err != nil {
			return err
		}

		str, ok := v.(string)
		if !ok {
			return errors.New("must be valid string")
		}
		s = str
	}

	if s == "" {
		v, err := validation.EnsureString(value)
		if err != nil {
			return err
		}
		s = v
	}

	if strings.Contains(s, " ") {
		return errors.New("must valid username with no spaces")
	}
	return nil
}
