package interestDTOs

import "time"

type Bio struct {
	InterestId string
	UserId     string
	Bio        string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Hobbie struct {
	Id     string
	Hobbie string
}

type MovieSerie struct {
	Id         string
	MovieSerie string
}

type Travel struct {
	Id     string
	Travel string
}

type Sport struct {
	Id    string
	Sport string
}

type InterestDetail struct {
	Bio
	Hobbies     []Hobbie
	MovieSeries []MovieSerie
	Travels     []Travel
	Sports      []Sport
}
