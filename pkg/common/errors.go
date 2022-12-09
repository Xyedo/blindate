package common

import (
	"errors"
	"net/http"
)

var (
	ErrUniqueConstraint23505 = &sentinelAPIError{status: http.StatusUnprocessableEntity, msg: "already created"}
	ErrRefNotFound23503      = &sentinelAPIError{status: http.StatusUnprocessableEntity, msg: "reference not found"}
	ErrNotMatchCredential    = &sentinelAPIError{status: http.StatusUnauthorized, msg: "invalid credentials"}
	ErrTooLongAccessingDB    = &sentinelAPIError{status: http.StatusConflict, msg: "request conflicted, please try again"}
	ErrResourceNotFound      = &sentinelAPIError{status: http.StatusNotFound, msg: "resource not found"}
)

var (
	ErrAuthorNotValid    = errors.New("author not in the conversation")
	ErrMaxProfilePicture = errors.New("excedeed profile picture constraint")
)
