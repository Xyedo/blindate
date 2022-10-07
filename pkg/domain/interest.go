package domain

import "time"

type Bio struct {
	Id        string    `json:"id,omitempty" db:"id"`
	UserId    string    `json:"userId" db:"user_id"`
	Bio       string    `json:"bio" db:"bio"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
type Interest struct {
	Bio         `json:"-" db:"-"`
	Hobbies     []Hobbie     `json:"hobbies" db:"-"`
	MovieSeries []MovieSerie `json:"movieSeries" db:"-"`
	Travels     []Travel     `json:"travels" db:"-"`
	Sports      []Sport      `json:"sports" db:"-"`
}
type Hobbie struct {
	Id     string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	Hobbie string `json:"hobbie" db:"hobbie" binding:"required,min=2,max=50"`
}
type MovieSerie struct {
	Id         string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	MovieSerie string `json:"movieSerie" db:"movie_serie" binding:"required,min=2,max=50"`
}
type Travel struct {
	Id     string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	Travel string `json:"travel" db:"travel" binding:"required,min=2,max=50"`
}
type Sport struct {
	Id    string `json:"id,omitempty" db:"id" binding:"omitempty,uuid"`
	Sport string `json:"sport" db:"sport" binding:"required,min=2,max=50"`
}
