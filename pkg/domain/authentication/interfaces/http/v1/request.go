package v1

import (
	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xyedo/blindate/pkg/common/mod"
)

type postAuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *postAuthRequest) Mod() *postAuthRequest {
	mod.Trim(&a.Email)
	return a
}
func (a postAuthRequest) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required),
	)
}
