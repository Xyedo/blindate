package dtos

import (
	"time"

	"github.com/xyedo/blindate/internal/optional"
)

type UpdateUser struct {
	Id          string                     `json:"-"`
	FullName    optional.Option[string]    `json:"fullName" binding:"omitempty,max=50"`
	Alias       optional.Option[string]    `json:"alias" binding:"omitempty,max=15"`
	Email       optional.Option[string]    `json:"email" binding:"omitempty,email"`
	OldPassword optional.Option[string]    `json:"oldPassword" binding:"required_with=NewPassword,omitempty,min=8"`
	NewPassword optional.Option[string]    `json:"newPassword" binding:"required_with=OldPassword,omitempty,min=8"`
	Dob         optional.Option[time.Time] `json:"dob" binding:"omitempty,validdob"`
}
