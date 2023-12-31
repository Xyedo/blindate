package entities

import apperror "github.com/xyedo/blindate/internal/common/app-error"

const (
	ErrCodeUserNotFound      apperror.Code = "USER_NOT_FOUND"
	ErrCodeInterestNotFound  apperror.Code = "INTEREST_NOT_FOUND"
	ErrCodeInterestTooLarge  apperror.Code = "INTEREST_TOO_LARGE"
	ErrCodeInterestDuplicate apperror.Code = "INTEREST_DUPLICATE"
	ErrCodePhotoInvalidType  apperror.Code = "PHOTO_INVALID_TYPE"
	ErrCodePhotoTooMuch      apperror.Code = "PHOTO_TOO_MUCH"
)
