package apperror

import "fmt"

type Sentinel struct {
	Err    error
	Message    string
	ErrMap map[string][]string
}

func (s Sentinel) Error() string {
	return s.Message
}


type Payload struct {
	Error   error
	Message string
}

func NotFound(payload Payload) error {
	if payload.Message == "" {
		payload.Message = "not found"
	}
	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrNotFound, payload.Error)
	} else {
		payload.Error = ErrNotFound
	}

	return Sentinel{
		Err: payload.Error,
		Message: payload.Message,
	}
}

func Forbidden(payload Payload) error {
	if payload.Message == "" {
		payload.Message = "you shouldn't access this resource"
	}

	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrForbiddenAccess, payload.Error)
	} else {
		payload.Error = ErrForbiddenAccess
	}
	return Sentinel{
		Err: payload.Error,
		Message: payload.Message,
	}
}

func Unauthorized(payload Payload) error {
	if payload.Message == "" {
		payload.Message = "unauthorized"
	}
	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrUnauthorized, payload.Error)
	} else {
		payload.Error = ErrUnauthorized
	}

	return Sentinel{
		Err: payload.Error,
		Message: payload.Message,
	}
}

func Conflicted(payload Payload) error {
	if payload.Message == "" {
		payload.Message = "resource conflicted, please try again!"
	}
	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrConflict, payload.Error)
	} else {
		payload.Error = ErrConflict
	}

	return Sentinel{
		Err: payload.Error,
		Message: payload.Message,
	}
}

func Timeout(payload Payload) error {
	if payload.Message == "" {
		payload.Message = "request timeout, please try again!"
	}
	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrTimeout, payload.Error)
	} else {
		payload.Error = ErrTimeout
	}

	return Sentinel{
		Err: payload.Error,
		Message: payload.Message,
	}
}

type PayloadMap struct {
	Error    error
	ErrorMap map[string][]string
}

func UnprocessableEntity(payload PayloadMap) error {
	if payload.ErrorMap == nil {
		payload.ErrorMap = map[string][]string{"unknown": {"unprocessable entity"}}
	}

	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrUnprocessableEntity, payload.Error)
	} else {
		payload.Error = ErrUnprocessableEntity
	}

	return Sentinel{
		Err:    payload.Error,
		ErrMap: payload.ErrorMap,
	}
}
