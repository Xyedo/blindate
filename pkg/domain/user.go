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
	Bio         `json:"-" db:"-"`
	Hobbies     []Hobbie     `json:"hobbies" db:"-"`
	MovieSeries []MovieSerie `json:"movieSeries" db:"-"`
	Travels     []Travel     `json:"travels" db:"-"`
	Sports      []Sport      `json:"sports" db:"-"`
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
