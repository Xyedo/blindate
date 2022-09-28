package domain

import (
	"errors"
)

var (
	ErrDuplicateEmail     = errors.New("database: duplicate email")
	ErrNotMatchCredential = errors.New("database: not match credential")
	ErrTooLongAccesingDB  = errors.New("database: too long accessing DB")
	ErrDuplicateToken     = errors.New("database: duplicate token")
	ErrResourceNotFound   = errors.New("database: resource not found")
)
