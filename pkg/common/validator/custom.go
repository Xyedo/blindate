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
	today := time.Now()
	oldest := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	if !(date.Before(today) && date.After(oldest)) {
		return errors.New("invalid day of birth")
	}
	return nil
}

func ValidUsername(value any) error {
	s, err := validation.EnsureString(value)
	if err != nil {
		if v, ok := value.(interface{ Value() (driver.Value, error) }); ok {
			v, err := v.Value()
			if err != nil {
				return err
			}
			if _, ok := v.(string); !ok {
				return errors.New("must be valid string")
			}
			return nil
		}
		return err
	}
	if strings.Contains(s, " ") {
		return errors.New("must valid username with no spaces")
	}
	return nil
}
