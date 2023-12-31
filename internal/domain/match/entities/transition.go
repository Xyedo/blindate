package entities

import (
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	userentities "github.com/xyedo/blindate/internal/domain/user/entities"
)

func (m Match) ValidateResource(requester userentities.UserDetail) error {
	if !(m.RequestFrom == requester.UserId || m.RequestTo == requester.UserId) {
		return apperror.Forbidden(apperror.Payload{
			Status: ErrCodeMatchIdInvalid,
		})
	}

	return nil
}
