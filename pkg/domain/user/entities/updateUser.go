package userEntity

import "time"

type Update struct {
	FullName    *string    `json:"fullName" binding:"omitempty,max=50"`
	Alias       *string    `json:"alias" binding:"omitempty,max=15"`
	Email       *string    `json:"email" binding:"omitempty,email"`
	OldPassword *string    `json:"oldPassword" binding:"required_with=NewPassword,omitempty,min=8"`
	NewPassword *string    `json:"newPassword" binding:"required_with=OldPassword,omitempty,min=8"`
	Dob         *time.Time `json:"dob" binding:"omitempty,validdob"`
}
