package userDTOs

import validation "github.com/go-ozzo/ozzo-validation"

type GetUserDetail struct {
	Id             string
	ProfilePicture bool
}

func (u GetUserDetail) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Id, validation.Required),
	)
}
