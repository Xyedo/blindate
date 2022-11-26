package domain

import (
	"errors"
)

var (
	ErrUniqueConstraint23505 = errors.New("database: violate unique constraint")
	ErrRefNotFound23503      = errors.New("database: reference not found")
	ErrNotMatchCredential    = errors.New("database: not match credential")
	ErrTooLongAccessingDB    = errors.New("database: too long accessing DB")
	ErrResourceNotFound      = errors.New("database: resource not found")
)
