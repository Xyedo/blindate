package entities

import (
	"time"

	"github.com/xyedo/blindate/pkg/optional"
)

type User struct {
	Id        string
	IsDeleted bool
}

type UserDetail struct {
	UserId           string
	Geog             Geography
	Bio              string
	LastOnline       time.Time
	Gender           Gender
	FromLoc          optional.String
	Height           optional.Int16
	EducationLevel   optional.String
	Drinking         optional.String
	Smoking          optional.String
	RelationshipPref optional.String
	LookingFor       optional.String
	Zodiac           optional.String
	Kids             optional.Int16
	Work             optional.String
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Version          int64

	Hobbies     []Hobbie
	MovieSeries []MovieSerie
	Travels     []Travel
	Sports      []Sport
}

type Hobbie struct {
	UUID      string
	UserId    string
	Hobbie    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64
}

type MovieSerie struct {
	UUID       string
	UserId     string
	MovieSerie string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Version    int64
}

type Travel struct {
	UUID      string
	UserId    string
	Travel    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64
}

type Sport struct {
	UUID      string
	UserId    string
	Sport     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64
}
