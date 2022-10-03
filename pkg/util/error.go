package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ReadValidationErr(err error, validation map[string]string) map[string]string {
	errMap := map[string]string{}
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		for _, err := range validationErr {
			if _, exist := errMap[err.Namespace()]; !exist {
				errMes, iexist := validation[err.Namespace()]
				if iexist {
					errMap[err.Namespace()] = errMes
				} else {
					errMap[err.Namespace()] = fmt.Sprintf("error validation on %s", err.Field())
				}

			}
		}
	}
	return errMap
}

func ReadJSONDecoderErr(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError
	switch {
	case errors.As(err, &syntaxError):
		return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("body contains badly-formed JSON")

	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

	case errors.Is(err, io.EOF):
		return errors.New("body must not be empty")

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Errorf("body contains unknown key %s", fieldName)

	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	default:
		return nil
	}
}
