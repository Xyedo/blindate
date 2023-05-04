package interestEntities

import "time"

type Bio struct {
	InterestId string    `db:"id"`
	UserId     string    `db:"user_id"`
	Bio        string    `db:"bio"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Hobbie struct {
	Id     string `db:"id"`
	Hobbie string `db:"hobbie"`
}

type MovieSerie struct {
	Id         string `db:"id"`
	MovieSerie string `db:"movie_serie"`
}

type Travel struct {
	Id     string `db:"id"`
	Travel string `db:"travel"`
}

type Sport struct {
	Id    string `db:"id"`
	Sport string `db:"sport"`
}
