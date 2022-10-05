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
	Id     string `json:"id,omitempty" db:"id"`
	Hobbie string `json:"hobbie" db:"hobbie"`
}
type MovieSerie struct {
	Id         string `json:"id,omitempty" db:"id"`
	MovieSerie string `json:"movieSerie" db:"movie_serie"`
}
type Travel struct {
	Id     string `json:"id,omitempty" db:"id"`
	Travel string `json:"travel" db:"travel"`
}
type Sport struct {
	Id    string `json:"id,omitempty" db:"id"`
	Sport string `json:"sport" db:"sport"`
}
