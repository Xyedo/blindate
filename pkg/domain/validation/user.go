package validation

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidDob(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	today := time.Now()
	oldest := time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)

	return (date.Before(today) && date.After(oldest))
}
