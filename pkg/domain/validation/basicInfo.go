package validation

import "github.com/go-playground/validator/v10"

func ValidEducationLevel(fl validator.FieldLevel) bool {
	educationLevelEnums := []string{"Less than high school diploma", "High school", "Some college, no degree", "Assosiate's Degree", "Bachelor's Degree", "Master's Degree", "Professional Degree", "Doctorate Degree"}
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	for _, validEnums := range educationLevelEnums {
		if validEnums == data {
			return true
		}
	}
	return false
}
