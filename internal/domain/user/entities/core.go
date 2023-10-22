package entities

import (
	"time"

	"github.com/xyedo/blindate/pkg/optional"
)

type User struct {
	Id        string
	IsDeleted bool
}
type BasicInfo struct {
	UserId           string
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
}
