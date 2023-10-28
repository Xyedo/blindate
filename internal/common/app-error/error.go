package apperror

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
	ErrForbiddenAccess     = errors.New("forbidden access")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrConflictIdempotent  = errors.New("conflict with idempotent")
	ErrConflict            = errors.New("conflict")
	ErrTimeout             = errors.New("timeout")
	ErrBadRequest          = errors.New("bad request")
)
