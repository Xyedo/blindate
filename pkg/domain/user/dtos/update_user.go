package userDTOs

import (
	"github.com/xyedo/blindate/internal/optional"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
)

type UpdateUser struct {
	Id          string
	FullName    optional.String
	Alias       optional.String
	Email       optional.String
	OldPassword optional.String
	NewPassword optional.String
	Dob         optional.Time
}

func (u UpdateUser) Validate() error {
	if !u.FullName.ValueSet() &&
		!u.Alias.ValueSet() &&
		!u.Email.ValueSet() &&
		!u.OldPassword.ValueSet() &&
		!u.NewPassword.ValueSet() &&
		!u.Dob.ValueSet() {
		return apperror.BadPayload(apperror.Payload{
			Message: "body must not be empty",
		})
	}
	return nil
}
