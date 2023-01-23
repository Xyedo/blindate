package apiError

import (
	"net/http"
)

var (
	ErrUniqueConstraint23505 = &sentinel{status: http.StatusUnprocessableEntity, msg: "already created"}
	ErrRefNotFound23503      = &sentinel{status: http.StatusUnprocessableEntity, msg: "reference not found"}
	ErrNotMatchCredential    = &sentinel{status: http.StatusUnauthorized, msg: "invalid credentials"}
	ErrTooLongAccessingDB    = &sentinel{status: http.StatusConflict, msg: "request conflicted, please try again"}
	ErrResourceNotFound      = &sentinel{status: http.StatusNotFound, msg: "resource not found"}
	ErrForbiddenAccess       = &sentinel{status: http.StatusForbidden, msg: "you shouldnt access this resource"}
)
