package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	latregex = regexp.MustCompile(`^[+-]?(([1-8]?[0-9])(\.[0-9]{1,6})?|90(\.0{1,6})?)$`)
	lngRegex = regexp.MustCompile(`^[+-]?((([1-9]?[0-9]|1[0-7][0-9])(\.[0-9]{1,6})?)|180(\.0{1,6})?)$`)
)

func ValidLat(fl validator.FieldLevel) bool {
	str, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return latregex.MatchString(str)

}
func ValidLng(fl validator.FieldLevel) bool {
	str, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return lngRegex.MatchString(str)
}
