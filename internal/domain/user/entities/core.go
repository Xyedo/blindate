package entities

import (
	"time"

	"github.com/xyedo/blindate/pkg/optional"
)

type User struct {
	Id        string
	IsDeleted bool
}

type UserDetails []UserDetail

type UserDetail struct {
	UserId           string
	Alias            string
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
	LookingFor       string
	Zodiac           optional.String
	Kids             optional.Int16
	Work             optional.String
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Version          int64

	Hobbies         []Hobbie
	MovieSeries     []MovieSerie
	Travels         []Travel
	Sports          []Sport
	ProfilePictures []ProfilePicture
}

type Hobbie struct {
	Id        string
	UserId    string
	Hobbie    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64
}

type MovieSerie struct {
	Id         string
	UserId     string
	MovieSerie string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Version    int64
}

type Travel struct {
	Id        string
	UserId    string
	Travel    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64
}

type Sport struct {
	Id        string
	UserId    string
	Sport     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int64
}

type ProfilePicture struct {
	Id           string
	UserId       string
	Selected     bool
	FileId       string
	presignedURL string
}

func (p *ProfilePicture) SetPresignedURL(url string) {
	p.presignedURL = url
}

func (p ProfilePicture) GetPresignedUrl() string {
	return p.presignedURL
}
