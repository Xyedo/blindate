package entities

import (
	"github.com/xyedo/blindate/pkg/optional"
)

type GetUserDetailOption struct {
	PessimisticLocking  bool
	WithHobbies         bool
	WithMovieSeries     bool
	WithTravels         bool
	WithSports          bool
	WithProfilePictures bool
}

type CreateUserDetail struct {
	Gender           string
	Geog             Geography
	Bio              string
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
}

type UpdateUserDetail struct {
	Gender           optional.String
	Geog             optional.Option[Geography]
	Bio              optional.String
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
}
