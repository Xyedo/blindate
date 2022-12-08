package common

import "errors"

type APIError interface {
	APIError() (int, string)
}
type sentinelAPIError struct {
	status int
	msg    string
}

func (e sentinelAPIError) Error() string {
	return e.msg
}
func (e sentinelAPIError) APIError() (int, string) {
	return e.status, e.msg
}

type sentinelWrappedError struct {
	error
	sentinel *sentinelAPIError
}

func (e sentinelWrappedError) Is(err error) bool {
	return errors.Is(err, e.sentinel)
}
func (e sentinelWrappedError) APIError() (int, string) {
	return e.sentinel.APIError()
}

func WrapError(err error, sentinel *sentinelAPIError) error {
	return sentinelWrappedError{error: err, sentinel: sentinel}
}
func WrapWithNewError(err error, status int, msg string) error {
	return sentinelWrappedError{error: err, sentinel: &sentinelAPIError{status: status, msg: msg}}
}

func WrapErrorWithMsg(err error, sentinel *sentinelAPIError, msg string) error {
	wrapedErr := sentinelWrappedError{error: err, sentinel: sentinel}
	wrapedErr.sentinel.msg = msg
	return wrapedErr
}
