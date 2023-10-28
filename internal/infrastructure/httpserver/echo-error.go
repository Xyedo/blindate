package httpserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/invopop/validation"
	"github.com/labstack/echo/v4"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
)

type Error struct {
	Message string                  `json:"message,omitempty"`
	Errors  []apperror.ErrorPayload `json:"errors,omitempty"`
}

func EchoErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	// application error
	var appErr apperror.Sentinel
	if errors.As(err, &appErr) {
		err := Error{
			Message: appErr.Message,
			Errors:  appErr.Payloads,
		}
		switch {
		case errors.Is(appErr.Err, apperror.ErrBadRequest):
			_ = c.JSON(http.StatusBadRequest, err)
			return

		case errors.Is(appErr.Err, apperror.ErrUnprocessableEntity):
			_ = c.JSON(http.StatusUnprocessableEntity, err)
			return

		case errors.Is(appErr.Err, apperror.ErrNotFound):
			_ = c.JSON(http.StatusNotFound, err)
			return

		case errors.Is(appErr.Err, apperror.ErrUnauthorized):
			_ = c.JSON(http.StatusUnauthorized, err)
			return

		case errors.Is(appErr.Err, apperror.ErrForbiddenAccess):
			_ = c.JSON(http.StatusForbidden, err)
			return

		case errors.Is(appErr.Err, apperror.ErrConflictIdempotent):
			_ = c.JSON(http.StatusOK, Error{
				Message: "already created",
			})
			return
		case errors.Is(appErr.Err, apperror.ErrConflict):
			_ = c.JSON(http.StatusConflict, err)
			return

		case errors.Is(appErr.Err, apperror.ErrTimeout):
			_ = c.JSON(http.StatusGatewayTimeout, err)
			return

		}
	}

	// json decoder error
	if field, msg := jsonDecoderError(err); msg != "" {
		if field != "" {
			_ = c.JSON(http.StatusBadRequest, Error{
				Message: "malformed_body",
				Errors: []apperror.ErrorPayload{
					{
						Status: apperror.StatusErrorMalformedRequestBody,
						Details: map[string][]any{
							field: {msg},
						},
					},
				},
			})
			return
		}
		_ = c.JSON(http.StatusBadRequest, Error{
			Message: msg,
		})
		return

	}

	var validatorError validation.Errors
	if errors.As(err, &validatorError) {
		mapErr := validationErrorMapping(validatorError)
		_ = c.JSON(http.StatusBadRequest, Error{
			Message: "validation_error",
			Errors: []apperror.ErrorPayload{
				{
					Status:  apperror.StatusErrorValidation,
					Details: mapErr,
				},
			},
		})
		return
	}

	//formfile error
	if errors.Is(err, http.ErrNotMultipart) || errors.Is(err, http.ErrMissingBoundary) {
		_ = c.JSON(http.StatusBadRequest, Error{
			Message: "content-type header is invalid",
		})
		return
	}
	if errors.Is(err, http.ErrMissingFile) {
		_ = c.JSON(http.StatusBadRequest, Error{
			Message: "request did not contain a file",
		})
		return
	}
	if errors.Is(err, multipart.ErrMessageTooLarge) {
		_ = c.JSON(http.StatusRequestEntityTooLarge, Error{
			Message: "file to large",
		})
		return
	}
	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		if echoErr.Internal != nil {
			c.Logger().Error(err)
		}
		_ = c.JSON(echoErr.Code, map[string]any{
			"message": echoErr.Message,
		})
		return

	}
	c.Echo().Logger.Error(err)

	_ = c.JSON(http.StatusInternalServerError, Error{
		Message: "cant catch 'em all, sorry!",
	})

}

func jsonDecoderError(err error) (field, message string) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError
	switch {
	case errors.As(err, &syntaxError):
		return "", fmt.Sprintf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

	case errors.Is(err, io.ErrUnexpectedEOF):
		return "", "body contains badly-formed JSON"

	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			var translatedType string
			switch unmarshalTypeError.Type.Name() {
			// REGEX *int*
			case "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "float32", "float64":
				translatedType = "number"
			case "Time":
				translatedType = "date time"
			case "string":
				translatedType = "string"
			}
			return unmarshalTypeError.Field, fmt.Sprintf("the field must be a valid %s", translatedType)
		}
		return "", fmt.Sprintf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

	case errors.Is(err, io.EOF):
		return "", "body must not be empty"

	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		fieldName, _ = strconv.Unquote(fieldName)
		return fieldName, "unknown key"
	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	default:
		return "", ""
	}
}
func validationErrorMapping(validatorError validation.Errors) map[string][]any {
	mapErr := make(map[string][]any)
	for key, err := range validatorError {
		if errs, ok := err.(validation.Errors); ok {
			newMap := validationErrorMapping(errs)
			mapErr = mergeMapWithKey(key, mapErr, newMap)
		} else {
			mapErr[key] = append(mapErr[key], err.Error())
		}

	}
	return mapErr
}

func mergeMapWithKey(key string, maps ...map[string][]any) map[string][]any {
	res := make(map[string][]any)
	for _, m := range maps {
		for k, v := range m {
			mergedKey := key + "." + k
			res[mergedKey] = append(res[mergedKey], v...)
		}
	}
	if len(res) == 0 {
		return nil
	}
	return res
}
