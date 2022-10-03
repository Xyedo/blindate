package entity

import "time"

type Interest struct {
	Id           string `db:"id"`
	UserId       string `db:"user_id"`
	Hobbies      []Hobies
	MoviesSeries []MovieSeries
	Traveling    []Traveling
	Sports       []Sports
	Bio          string    `db:"bio"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
type Hobies struct {
	Id     string `db:"id"`
	Hobbie string `db:"hobbie"`
}
type MovieSeries struct {
	Id         string `db:"id"`
	MovieSerie string `db:"movie_serie"`
}
type Traveling struct {
	Id     string `db:"id"`
	Travel string `db:"travel"`
}
type Sports struct {
	Id    string `db:"id"`
	Sport string `db:"sport"`
}
