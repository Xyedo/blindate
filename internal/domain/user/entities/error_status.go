package entities

import apperror "github.com/xyedo/blindate/internal/common/app-error"

const (
	UserNotFound      apperror.StatusError = "USER_NOT_FOUND"
	InterestNotFound  apperror.StatusError = "INTEREST_NOT_FOUND"
	InterestTooLarge  apperror.StatusError = "INTEREST_TOO_LARGE"
	InterestDuplicate apperror.StatusError = "INTEREST_DUPLICATE"
	PhotoInvalidType  apperror.StatusError = "PHOTO_INVALID_TYPE"
	PhotoTooMuch      apperror.StatusError = "PHOTO_TOO_MUCH"
)
