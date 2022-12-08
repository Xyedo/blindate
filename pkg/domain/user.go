package domain

import (
	"time"
)

type User struct {
	ID             string           `db:"id" json:"id"`
	FullName       string           `db:"full_name" json:"fullName"`
	Alias          string           `db:"alias" json:"alias"`
	Email          string           `db:"email" json:"email"`
	ProfilePic     []ProfilePicture `db:"-" json:"profilePicture,omitempty"`
	Password       string           `db:"-" json:"-"`
	HashedPassword string           `db:"password" json:"-"`
	Active         bool             `db:"active" json:"-"`
	Dob            time.Time        `db:"dob" json:"dob"`
	CreatedAt      time.Time        `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time        `db:"updated_at" json:"updatedAt"`
}

type UpdateUser struct {
	FullName    *string    `json:"fullName" binding:"omitempty,max=50"`
	Alias       *string    `json:"alias" binding:"omitempty,max=15"`
	Email       *string    `json:"email" binding:"omitempty,email"`
	OldPassword *string    `json:"oldPassword" binding:"required_with=NewPassword,omitempty,min=8"`
	NewPassword *string    `json:"newPassword" binding:"required_with=OldPassword,omitempty,min=8"`
	Dob         *time.Time `json:"dob" binding:"omitempty,validdob"`
}

// Bio one to one with user
type Bio struct {
	Id        string    `json:"id,omitempty" db:"id"`
	UserId    string    `json:"userId" db:"user_id"`
	Bio       string    `json:"bio" db:"bio"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// BasicInfo one to one with user
type BasicInfo struct {
	UserId           string    `json:"userId"`
	Gender           string    `json:"gender"`
	FromLoc          *string   `json:"fromLoc"`
	Height           *int      `json:"height"`
	EducationLevel   *string   `json:"educationLevel"`
	Drinking         *string   `json:"drinking"`
	Smoking          *string   `json:"smoking"`
	RelationshipPref *string   `json:"relationshipPref"`
	LookingFor       string    `json:"lookingFor"`
	Zodiac           *string   `json:"zodiac"`
	Kids             *int      `json:"kids"`
	Work             *string   `json:"work"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
type UpdateBasicInfo struct {
	Gender           *string `json:"gender" binding:"omitempty,max=25"`
	FromLoc          *string `json:"fromLoc" binding:"omitempty,max=25"`
	Height           *int    `json:"height" binding:"omitempty,min=0,max=300"`
	EducationLevel   *string `json:"educationLevel" binding:"omitempty,max=49"`
	Drinking         *string `json:"drinking" binding:"omitempty,max=49"`
	Smoking          *string `json:"smoking" binding:"omitempty,max=49"`
	RelationshipPref *string `json:"relationshipPref" binding:"omitempty,max=49"`
	LookingFor       *string `json:"lookingFor"  binding:"omitempty,max=25"`
	Zodiac           *string `json:"zodiac" binding:"omitempty,max=50"`
	Kids             *int    `json:"kids" binding:"omitempty,max=100"`
	Work             *string `json:"work" binding:"omitempty,max=50"`
}

// Location one to one with user
type Location struct {
	UserId string `json:"-"`
	Lat    string `json:"lat"`
	Lng    string `json:"lng"`
}

// Online one to one with user
type Online struct {
	UserId     string    `json:"-" db:"user_id"`
	LastOnline time.Time `json:"lastOnline" db:"last_online"`
	IsOnline   bool      `json:"isOnline" db:"is_online"`
}

// Interest one to one with user
type Interest struct {
	Bio
	Hobbies     []Hobbie     `json:"hobbies"`
	MovieSeries []MovieSerie `json:"movieSeries"`
	Travels     []Travel     `json:"travels"`
	Sports      []Sport      `json:"sports"`
}

// ProfilePicture one to many with user
type ProfilePicture struct {
	Id          string `json:"id" db:"id"`
	UserId      string `json:"userId" db:"user_id"`
	Selected    bool   `json:"selected" db:"selected"`
	PictureLink string `json:"pictureLink" db:"picture_ref"`
}

// Hobbie one to many with bio
type Hobbie struct {
	Id     string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	Hobbie string `json:"hobbie" db:"hobbie" binding:"required,min=2,max=50"`
}

// MovieSerie one to many with bio
type MovieSerie struct {
	Id         string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	MovieSerie string `json:"movieSerie" db:"movie_serie" binding:"required,min=2,max=50"`
}

// Travel one to many with bio
type Travel struct {
	Id     string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	Travel string `json:"travel" db:"travel" binding:"required,min=2,max=50"`
}

// Sport one to many with bio
type Sport struct {
	Id    string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	Sport string `json:"sport" db:"sport" binding:"required,min=2,max=50"`
}
