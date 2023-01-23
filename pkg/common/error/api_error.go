package apiError

import "errors"

type API interface {
	APIError() (int, string)
}
type sentinel struct {
	status int
	msg    string
}

func (e sentinel) Error() string {
	return e.msg
}
func (e sentinel) APIError() (int, string) {
	return e.status, e.msg
}

type sentinelWrappedError struct {
	error
	sentinel *sentinel
}

func (e sentinelWrappedError) Is(err error) bool {
	return errors.Is(err, e.sentinel)
}
func (e sentinelWrappedError) APIError() (int, string) {
	return e.sentinel.APIError()
}

func Wrap(err error, sentinel *sentinel) error {
	return sentinelWrappedError{error: err, sentinel: sentinel}
}
func WrapWithNewSentinel(err error, status int, msg string) error {
	return sentinelWrappedError{error: err, sentinel: &sentinel{status: status, msg: msg}}
}

func WrapWithMsg(err error, sentinel *sentinel, msg string) error {
	wrapedErr := sentinelWrappedError{error: err, sentinel: sentinel}
	wrapedErr.sentinel.msg = msg
	return wrapedErr
}
