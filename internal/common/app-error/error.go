package apperror

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrForbiddenAccess     = errors.New("forbidden access")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrConflict            = errors.New("conflict")
	ErrTimeout             = errors.New("timeout")
	ErrBadRequst           = errors.New("bad request")
)
