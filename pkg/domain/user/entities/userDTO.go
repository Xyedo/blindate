package userEntity

import "time"

type FullDTO struct {
	ID         string       `db:"id" json:"id"`
	FullName   string       `db:"full_name" json:"fullName"`
	Alias      string       `db:"alias" json:"alias"`
	Email      string       `db:"email" json:"email"`
	ProfilePic []ProfilePic `db:"-" json:"profilePicture,omitempty"`
	Password   string       `db:"password" json:"-"`
	Active     bool         `db:"active" json:"-"`
	Dob        time.Time    `db:"dob" json:"dob"`
	CreatedAt  time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time    `db:"updated_at" json:"updatedAt"`
}
