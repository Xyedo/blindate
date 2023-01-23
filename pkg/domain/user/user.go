package user

import "errors"

var (
	ErrNoUserToMatch     = errors.New("no user to match")
	ErrMaxProfilePicture = errors.New("excedeed profile picture constraint")
)
