package v1

import (
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xyedo/blindate/internal/optional"
	"github.com/xyedo/blindate/pkg/common/mod"
	"github.com/xyedo/blindate/pkg/common/validator"
)

type postUserRequest struct {
	FullName string    `json:"full_name" mold:"trim_whitespace"`
	Alias    string    `json:"alias" mod:"trim"`
	Email    string    `json:"email" mod:"trim"`
	Password string    `json:"password"`
	Dob      time.Time `json:"dob"`
}

func (u *postUserRequest) mod() *postUserRequest {
	mod.TrimWhiteSpace(&u.FullName)
	mod.Trim(&u.Alias)
	mod.Trim(&u.Email)
	return u
}

func (u postUserRequest) validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FullName, validation.Required, validation.Length(1, 50)),
		validation.Field(&u.Alias, validation.Required, validation.By(validator.ValidUsername), validation.Length(1, 15)),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 0)),
		validation.Field(&u.Dob, validation.Required, validation.By(validator.ValidDob)),
	)
}

type patchUserRequest struct {
	FullName    optional.String `json:"full_name"  mold:"trim_whitespace"`
	Alias       optional.String `json:"alias" mod:"trim"`
	Email       optional.String `json:"email" mod:"trim"`
	OldPassword optional.String `json:"old_password"`
	NewPassword optional.String `json:"new_password"`
	Dob         optional.Time   `json:"dob"`
}

func (u *patchUserRequest) mod() *patchUserRequest {
	u.FullName.If(func(s string) {
		mod.TrimWhiteSpace(&s)
		u.FullName.Set(s)
	})
	u.Alias.If(func(s string) {
		mod.Trim(&s)
		u.Alias.Set(s)
	})

	u.Email.If(func(s string) {
		mod.Trim(&s)
		u.Email.Set(s)
	})
	return u
}

func (u patchUserRequest) validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FullName, validation.Skip.When(!u.FullName.ValueSet()), validation.Length(1, 50)),
		validation.Field(&u.Alias, validation.Skip.When(!u.Alias.ValueSet()), validation.By(validator.ValidUsername), validation.Length(1, 15)),
		validation.Field(&u.Email, validation.Skip.When(!u.Email.ValueSet()), is.Email),
		validation.Field(&u.OldPassword, validation.Required.When(u.NewPassword.Present()), validation.Skip.When(!u.OldPassword.ValueSet()), validation.Length(8, 0)),
		validation.Field(&u.NewPassword, validation.Required.When(u.OldPassword.Present()), validation.Skip.When(!u.NewPassword.ValueSet()), validation.Length(8, 0)),
		validation.Field(&u.Dob, validation.Skip.When(!u.Dob.ValueSet()), validation.By(validator.ValidDob)),
	)
}
