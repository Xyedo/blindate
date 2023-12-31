package entities

import apperror "github.com/xyedo/blindate/internal/common/app-error"

const (
	ErrCodeMatchCandidateEmpty apperror.Code = "MATCH_CANDIDATE_EMPTY"
	ErrCodeMatchNotFound       apperror.Code = "MATCH_NOT_FOUND"
	ErrCodeMatchIdInvalid      apperror.Code = "MATCH_ID_INVALID"
	ErrCodeMatchStatusInvalid  apperror.Code = "MATCH_STATUS_INVALID"
)
