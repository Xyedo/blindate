package entities

import apperror "github.com/xyedo/blindate/internal/common/app-error"

const (
	UserNotFound      apperror.Code = "USER_NOT_FOUND"
	InterestNotFound  apperror.Code = "INTEREST_NOT_FOUND"
	InterestTooLarge  apperror.Code = "INTEREST_TOO_LARGE"
	InterestDuplicate apperror.Code = "INTEREST_DUPLICATE"
	PhotoInvalidType  apperror.Code = "PHOTO_INVALID_TYPE"
	PhotoTooMuch      apperror.Code = "PHOTO_TOO_MUCH"
)
