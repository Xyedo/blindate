package apperror

import "fmt"

type StatusError string

const (
	StatusErrorMalformedRequestBody StatusError = "MALFORMED_REQUEST_BODY"
	StatusErrorValidation           StatusError = "VALIDATION_ERROR"
	StatusErrorInvalidAuth          StatusError = "INVALID_AUTHORIZATION"
)
const (
	statusErrorDefaultNotFound             StatusError = "NOT_FOUND"
	statusErrorDefaultForbidden            StatusError = "FORBIDDEN"
	statusErrorDefaultUnauthorized         StatusError = "UNAUTHORIZED"
	statusErrorDefaultConflicted           StatusError = "CONFLICTED"
	statusErrorDefaultTimeout              StatusError = "TIMEOUT"
	statusErrorDefaultBadPayload           StatusError = "BAD_PAYLOAD"
	statusErrorDefaultUnprocessablePayload StatusError = "UNPROCESSABLE_PAYLOAD"
)

type Sentinel struct {
	Err      error
	Message  string
	Payloads []ErrorPayload
}

type ErrorPayload struct {
	Status StatusError
	ErrMap map[string][]any
}

func (s Sentinel) Error() string {
	return s.Err.Error()
}

type Payload struct {
	Error   error
	Status  StatusError
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
	if payload.Status == "" {
		payload.Status = statusErrorDefaultNotFound
	}
	return Sentinel{
		Err: payload.Error,
		Payloads: []ErrorPayload{
			{
				Status: payload.Status,
			},
		},
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

	if payload.Status == "" {
		payload.Status = statusErrorDefaultForbidden
	}
	return Sentinel{
		Err: payload.Error,
		Payloads: []ErrorPayload{
			{
				Status: payload.Status,
			},
		},
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

	if payload.Status == "" {
		payload.Status = statusErrorDefaultUnauthorized
	}

	return Sentinel{
		Err: payload.Error,
		Payloads: []ErrorPayload{
			{
				Status: payload.Status,
			},
		},
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

	if payload.Status == "" {
		payload.Status = statusErrorDefaultConflicted
	}

	return Sentinel{
		Err: payload.Error,
		Payloads: []ErrorPayload{
			{
				Status: payload.Status,
			},
		},
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

	if payload.Status == "" {
		payload.Status = statusErrorDefaultTimeout
	}

	return Sentinel{
		Err: payload.Error,
		Payloads: []ErrorPayload{
			{
				Status: payload.Status,
			},
		},
		Message: payload.Message,
	}
}

func BadPayload(payload Payload) error {
	if payload.Message == "" {
		payload.Message = "bad request"
	}

	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrBadRequest, payload.Error)
	} else {
		payload.Error = ErrBadRequest
	}

	if payload.Status == "" {
		payload.Status = statusErrorDefaultBadPayload
	}

	return Sentinel{
		Err: payload.Error,
		Payloads: []ErrorPayload{
			{
				Status: payload.Status,
			},
		},
		Message: payload.Message,
	}
}

type PayloadMap struct {
	Error    error
	Message  string
	Payloads []ErrorPayload
}

func BadPayloadWithPayloadMap(payload PayloadMap) error {
	if payload.Message == "" {
		payload.Message = "bad request"
	}

	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrBadRequest, payload.Error)
	} else {
		payload.Error = ErrBadRequest
	}

	if payload.Payloads == nil {
		payload.Payloads = []ErrorPayload{
			{
				Status: statusErrorDefaultBadPayload,
			},
		}
	} else {
		for i := range payload.Payloads {
			if payload.Payloads[i].Status == "" {
				payload.Payloads[i].Status = statusErrorDefaultBadPayload
			}
		}
	}

	return Sentinel{
		Err:      payload.Error,
		Payloads: payload.Payloads,
		Message:  payload.Message,
	}
}

func UnprocessableEntity(payload Payload) error {
	if payload.Message == "" {
		payload.Message = "we undestand your request, but its unprocessable"
	}

	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrUnprocessableEntity, payload.Error)
	} else {
		payload.Error = ErrUnprocessableEntity
	}

	if payload.Status == "" {
		payload.Status = statusErrorDefaultUnprocessablePayload
	}

	return Sentinel{
		Err:     payload.Error,
		Message: payload.Message,
		Payloads: []ErrorPayload{
			{
				Status: payload.Status,
			},
		},
	}
}

func UnprocessableEntityWithPayloadMap(payload PayloadMap) error {
	if payload.Message == "" {
		payload.Message = "we undestand your request, but its unprocessable"
	}

	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrUnprocessableEntity, payload.Error)
	} else {
		payload.Error = ErrUnprocessableEntity
	}

	if payload.Payloads == nil {
		payload.Payloads = []ErrorPayload{
			{
				Status: statusErrorDefaultUnprocessablePayload,
			},
		}
	} else {
		for i := range payload.Payloads {
			if payload.Payloads[i].Status == "" {
				payload.Payloads[i].Status = statusErrorDefaultUnprocessablePayload
			}
		}
	}

	return Sentinel{
		Err:      payload.Error,
		Message:  payload.Message,
		Payloads: payload.Payloads,
	}
}
