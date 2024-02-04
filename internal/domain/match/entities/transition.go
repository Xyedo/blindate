package entities

import (
	apperror "github.com/xyedo/blindate/internal/common/app-error"
)

func (m Match) ValidateResource(requestId string) (string, error) {
	if !(m.RequestFrom == requestId || m.RequestTo == requestId) {
		return "", apperror.Forbidden(apperror.Payload{
			Status: ErrCodeMatchIdInvalid,
		})
	}

	if m.RequestFrom == requestId {
		return m.RequestTo, nil
	}

	return m.RequestFrom, nil

}
