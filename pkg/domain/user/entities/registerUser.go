package userEntity

import "time"

type Register struct {
	FullName string    `json:"fullName" binding:"required,max=50"`
	Alias    string    `json:"alias" binding:"required,max=15"`
	Email    string    `json:"email" binding:"required,email"`
	Password string    `json:"password" binding:"required,min=8"`
	Dob      time.Time `json:"dob" binding:"required,validdob"`
}
