package validator

import (
	"database/sql/driver"
	"errors"
	"strings"
	"time"

	"github.com/invopop/validation"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
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
		return validation.NewError("validation_day_of_birth", "invalid day of birth")
	}
	return nil
}

func ValidUsername(value any) error {
	var s string
	v, ok := value.(interface{ Value() (driver.Value, error) })
	if ok {
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
		return validation.NewError("validation_valid_username", "value must be valid username with no spaces")
	}
	return nil
}

func Unique[T comparable](vs []T, key string) error {
	uniqueValue := make(map[T]bool)
	for _, v := range vs {
		if _, ok := uniqueValue[v]; ok {
			return apperror.UnprocessableEntity(apperror.PayloadMap{
				ErrorMap: map[string]string{
					key: "every value must be unique",
				},
			})
		}
		uniqueValue[v] = true
	}
	return nil
}
