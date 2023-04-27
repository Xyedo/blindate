package httperror

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
)

var regexSnakeCase = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

type Error struct {
	Message string              `json:"message,omitempty"`
	Errors  map[string][]string `json:"errors,omitempty"`
}

func HandleError(c *gin.Context, err error) {
	// application error
	var appErr *apperror.Sentinel
	if errors.As(err, &appErr) {
		switch {
		case errors.Is(err, apperror.ErrUnprocessableEntity):
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, Error{
				Message: appErr.Message,
				Errors:  appErr.ErrMap,
			})
			return
		case errors.Is(err, apperror.ErrNotFound):
			c.AbortWithStatus(http.StatusNotFound)
			return
		case errors.Is(err, apperror.ErrUnauthorized):
			c.AbortWithStatusJSON(http.StatusUnauthorized, Error{
				Message: appErr.Message,
			})
			return
		case errors.Is(err, apperror.ErrForbiddenAccess):
			c.AbortWithStatusJSON(http.StatusForbidden, Error{
				Message: appErr.Message,
			})
			return
		case errors.Is(err, apperror.ErrConflict):
			c.AbortWithStatusJSON(http.StatusConflict, Error{
				Message: appErr.Message,
			})
			return
		case errors.Is(err, apperror.ErrTimeout):
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, Error{
				Message: appErr.Message,
			})
			return
		}
	}

	// json decoder error
	if err := jsonDecoderError(err); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, Error{
			Message: err.Error(),
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

		case gin.ErrorTypeBind:
			errors := ginErr.Err.(validator.ValidationErrors)
			validationErrors := make(map[string][]string)
			for _, err := range errors {
				switch err.Tag() {
				case "required":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], "the field is required")
				case "min":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("the field must be at least %s", err.Param()))
				case "max":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("the field must not be greater than %s", err.Param()))
				case "numeric":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], "the field must be a valid number")
				case "required_with":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field is required when %s is present.", err.Field(), toSnakeCase(err.Param())))
				case "uuid4":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be a valid UUID.", err.Field()))
				case "gte":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be greater than or equal %s", err.Field(), err.Param()))
				case "gt":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be greater than %s", err.Field(), err.Param()))
				case "lte":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be less than or equal %s", err.Field(), err.Param()))
				case "lt":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be less than %s", err.Field(), err.Param()))
				case "email":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be a valid email address", err.Field()))
				case "iso8601":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be a valid ISO8601 time duration format", err.Field()))
				case "gtefield":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be greater than or equal the %s field", err.Field(), toSnakeCase(err.Param())))
				case "nefield":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be not equal to the %s field", err.Field(), toSnakeCase(err.Param())))
				case "gtetoday":
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The %s field must be a greater than or equal to today's date", err.Field()))
				default:
					validationErrors[err.Field()] = append(validationErrors[err.Field()], fmt.Sprintf("The error:%s was not registered on field %s", err.Tag(), err.Field()))

				}
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, Error{
				Message: "validation_error",
				Errors:  validationErrors,
			})
			return
		}

	}

	log.Println(err)
	c.AbortWithStatus(http.StatusInternalServerError)
}

func toSnakeCase(s string) string {
	return strings.ToLower(strings.Join(regexSnakeCase.FindAllString(s, -1), "_"))
}

func jsonDecoderError(err error) error {
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
