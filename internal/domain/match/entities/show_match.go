package entities

import apperror "github.com/xyedo/blindate/internal/common/app-error"

func (m Match) ValidateShow(requestId string) error {
	if m.RequestStatus == MatchStatusUnknown ||
		m.RequestStatus == MatchStatusDeclined ||
		m.RequestStatus == MatchStatusRequested && m.UpdatedBy.MustGet() == requestId {
		return apperror.Forbidden(apperror.Payload{
			Status: ErrCodeMatchIdInvalid,
		})
	}

	return nil
}
