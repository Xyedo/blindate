package httperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xyedo/blindate/internal/security"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
)

type Error struct {
	Message string              `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func HandleError(c *gin.Context, err error) {
	// application error
	var appErr apperror.Sentinel
	if errors.As(err, &appErr) {
		switch {
		case errors.Is(appErr.Err, apperror.ErrBadRequst):
			c.AbortWithStatusJSON(http.StatusBadRequest, Error{
				Message: appErr.Message,
			})
			return
		case errors.Is(appErr.Err, apperror.ErrUnprocessableEntity):
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, Error{
				Message: appErr.Message,
				Errors:  appErr.ErrMap,
			})
			return
		case errors.Is(appErr.Err, apperror.ErrNotFound):
			c.AbortWithStatus(http.StatusNotFound)
			return
		case errors.Is(appErr.Err, apperror.ErrUnauthorized):
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error{
				Message: appErr.Message,
			})
			return
		case errors.Is(appErr.Err, apperror.ErrForbiddenAccess):
			c.AbortWithStatusJSON(http.StatusForbidden, Error{
				Message: appErr.Message,
			})
			return
		case errors.Is(appErr.Err, apperror.ErrConflict):
			c.AbortWithStatusJSON(http.StatusConflict, Error{
				Message: appErr.Message,
			})
			return
		case errors.Is(appErr.Err, apperror.ErrTimeout):
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, Error{
				Message: appErr.Message,
			})
			return
		}
	}
	//jwtError
	if errors.Is(err, security.ErrInvalidCred) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Error{
			Message: "invalid credentials",
		})
		return
	}

	if errors.Is(err, security.ErrInvalidCred) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, Error{
			Message: "token expired",
		})
		return
	}

	// json decoder error
	if field, msg := jsonDecoderError(err); msg != "" {
		if field != "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, Error{
				Message: "decode json error",
				Errors:  map[string][]string{field: {msg}},
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, Error{
			Message: msg,
		})
		return
	}

	//usecase validation error
	var validatorError validation.Errors
	if errors.As(err, &validatorError) {
		c.AbortWithStatusJSON(http.StatusBadRequest, validatorError)
		return
	}

	//time parse error
	var timeParseErr *time.ParseError
	if errors.As(err, &timeParseErr) {
		c.AbortWithStatusJSON(http.StatusBadRequest, Error{
			Message: fmt.Sprintf("invalid time format on %s", timeParseErr.Value),
		})
		return
	}

	//formfile error
	if errors.Is(err, http.ErrNotMultipart) || errors.Is(err, http.ErrMissingBoundary) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "content-Type header is not valid"})
		return
	}
	if errors.Is(err, http.ErrMissingFile) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "request did not contain a file"})
		return
	}
	if errors.Is(err, multipart.ErrMessageTooLarge) {
		c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
			"message": "max byte to upload is 8mB",
		})
		return
	}

	//gin error
	var ginErr *gin.Error
	if errors.As(err, &ginErr) {
		switch ginErr.Type {
		case gin.ErrorTypePublic:
			if !c.Writer.Written() {
				c.JSON(c.Writer.Status(), gin.H{"message": ginErr.Err})
				return
			}
		}
		return
	}

	log.Println(err)
	log.Println(reflect.TypeOf(err))
	if !c.Writer.Written() {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
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
