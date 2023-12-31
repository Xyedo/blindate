package apperror

import "fmt"

type Code string

const (
	StatusErrorMalformedRequestBody Code = "MALFORMED_REQUEST_BODY"
	StatusErrorValidation           Code = "VALIDATION_ERROR"
	StatusErrorInvalidAuth          Code = "INVALID_AUTHORIZATION"
)
const (
	statusErrorDefaultDuplicate            Code = "DUPLICATE"
	statusErrorDefaultNotFound             Code = "NOT_FOUND"
	statusErrorDefaultForbidden            Code = "FORBIDDEN"
	statusErrorDefaultUnauthorized         Code = "UNAUTHORIZED"
	statusErrorDefaultConflicted           Code = "CONFLICTED"
	statusErrorDefaultTimeout              Code = "TIMEOUT"
	statusErrorDefaultBadPayload           Code = "BAD_PAYLOAD"
	statusErrorDefaultUnprocessablePayload Code = "UNPROCESSABLE_PAYLOAD"
)

type Sentinel struct {
	Err      error          `json:"-"`
	Payloads []ErrorPayload `json:"payload"`
}

type ErrorPayload struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func (s Sentinel) Error() string {
	return s.Err.Error()
}

type Payload struct {
	Error   error
	Status  Code
	Message string
}

func New(payload Payload, details ...any) error {
	return Sentinel{
		Err: payload.Error,
		Payloads: []ErrorPayload{
			{
				Code:    payload.Status,
				Message: payload.Message,
				Details: details,
			},
		},
	}
}

func Duplicate(payload Payload, indempotent bool) error {
	if payload.Message == "" {
		payload.Message = "duplicate"
	}

	var payloadErr error
	if indempotent {
		payloadErr = ErrConflictIdempotent
	} else {
		payloadErr = ErrConflict
	}
	if payload.Error != nil {
		payloadErr = fmt.Errorf("%w:%w", payloadErr, payload.Error)
	}

	if payload.Status == "" {
		payload.Status = statusErrorDefaultDuplicate
	}
	return Sentinel{
		Err: payloadErr,
		Payloads: []ErrorPayload{
			{
				Code:    payload.Status,
				Message: payload.Message,
			},
		},
	}
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
				Code:    payload.Status,
				Message: payload.Message,
			},
		},
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
				Code:    payload.Status,
				Message: payload.Message,
			},
		},
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
				Code:    payload.Status,
				Message: payload.Message,
			},
		},
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
				Code:    payload.Status,
				Message: payload.Message,
			},
		},
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
				Code:    payload.Status,
				Message: payload.Message,
			},
		},
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
				Code:    payload.Status,
				Message: payload.Message,
			},
		},
	}
}

type PayloadMap struct {
	Error    error
	Payloads []ErrorPayload
}

func BadPayloadWithPayloadMap(payload PayloadMap) error {

	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrBadRequest, payload.Error)
	} else {
		payload.Error = ErrBadRequest
	}

	if payload.Payloads == nil {
		payload.Payloads = []ErrorPayload{
			{
				Code:    statusErrorDefaultBadPayload,
				Message: "bad request",
			},
		}
	} else {
		for i := range payload.Payloads {
			if payload.Payloads[i].Code == "" {
				payload.Payloads[i].Code = statusErrorDefaultBadPayload
			}
		}
	}

	return Sentinel{
		Err:      payload.Error,
		Payloads: payload.Payloads,
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
		Err: payload.Error,
		Payloads: []ErrorPayload{
			{
				Message: payload.Message,
				Code:    payload.Status,
			},
		},
	}
}

func UnprocessableEntityWithPayloadMap(payload PayloadMap) error {
	if payload.Error != nil {
		payload.Error = fmt.Errorf("%w:%w", ErrUnprocessableEntity, payload.Error)
	} else {
		payload.Error = ErrUnprocessableEntity
	}

	if payload.Payloads == nil {
		payload.Payloads = []ErrorPayload{
			{
				Message: "we understand your request, but its unprocessable",
				Code:    statusErrorDefaultUnprocessablePayload,
			},
		}
	} else {
		for i := range payload.Payloads {
			if payload.Payloads[i].Code == "" {
				payload.Payloads[i].Code = statusErrorDefaultUnprocessablePayload
			}
		}
	}

	return Sentinel{
		Err:      payload.Error,
		Payloads: payload.Payloads,
	}
}
