package domain

import "errors"

var (
	ErrDuplicateEmail    = errors.New("database: duplicate email")
	ErrTooLongAccesingDB = errors.New("database: too long accessing DB")
)
