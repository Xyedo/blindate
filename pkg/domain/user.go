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
	Gender           string    `json:"gender" binding:"required,oneof=Female Male Other"`
	FromLoc          *string   `json:"fromLoc" binding:"omitempty,max=99"`
	Height           *int      `json:"height" binding:"omitempty,min=0,max=400"`
	EducationLevel   *string   `json:"educationLevel" binding:"omitempty,valideducationlevel"`
	Drinking         *string   `json:"drinking" binding:"omitempty,oneof=Never Ocassionally 'Once a week' 'More than 2/3 times a week' 'Every day'"`
	Smoking          *string   `json:"smoking" binding:"omitempty,oneof=Never Ocassionally 'Once a week' 'More than 2/3 times a week' 'Every day'"`
	RelationshipPref *string   `json:"relationshipPref" binding:"omitempty,oneof='One night Stand' 'Having fun' Serious"`
	LookingFor       string    `json:"lookingFor" binding:"required,oneof=Female Male Other"`
	Zodiac           *string   `json:"zodiac" binding:"omitempty,oneof=Aries Taurus Gemini Cancer Leo Virgo Libra Scorpio Sagittarius Capricorn Aquarius Pisces"`
	Kids             *int      `json:"kids" binding:"omitempty,min=0,max=30"`
	Work             *string   `json:"work" binding:"omitempty,max=50"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
type UpdateBasicInfo struct {
	Gender           *string `json:"gender" binding:"omitempty,oneof=Female Male Other"`
	FromLoc          *string `json:"fromLoc" binding:"omitempty,max=99"`
	Height           *int    `json:"height" binding:"omitempty,min=0,max=400"`
	EducationLevel   *string `json:"educationLevel" binding:"omitempty,valideducationlevel"`
	Drinking         *string `json:"drinking" binding:"omitempty,oneof=Never Ocassionally 'Once a week' 'More than 2/3 times a week' 'Every day'"`
	Smoking          *string `json:"smoking" binding:"omitempty,oneof=Never Ocassionally 'Once a week' 'More than 2/3 times a week' 'Every day'"`
	RelationshipPref *string `json:"relationshipPref" binding:"omitempty,oneof='One night Stand' 'Having fun' Serious"`
	LookingFor       *string `json:"lookingFor"  binding:"omitempty,oneof=Female Male Other"`
	Zodiac           *string `json:"zodiac" binding:"omitempty,oneof=Aries Taurus Gemini Cancer Leo Virgo Libra Scorpio Sagittarius Capricorn Aquarius Pisces"`
	Kids             *int    `json:"kids" binding:"omitempty,min=0,max=30"`
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
