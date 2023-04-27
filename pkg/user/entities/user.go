package entities

import "time"

type UserProfilePic struct {
	Id          string `db:"id"`
	UserId      string `db:"user_id"`
	Selected    bool   `db:"selected"`
	PictureLink string `db:"picture_ref"`
}

type User struct {
	ID         string           `db:"id"`
	FullName   string           `db:"full_name"`
	Alias      string           `db:"alias"`
	ProfilePic []UserProfilePic `db:"-"`
	Email      string           `db:"email"`
	Password   string           `db:"password"`
	Dob        time.Time        `db:"dob"`
	Active     bool             `db:"active"`
	CreatedAt  time.Time        `db:"created_at"`
	UpdatedAt  time.Time        `db:"updated_at"`
}
