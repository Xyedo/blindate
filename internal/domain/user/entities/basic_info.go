package entities

import (
	"time"

	"github.com/xyedo/blindate/pkg/optional"
)

type CreateBasicInfo struct {
	Gender           string
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

type UpdateBasicInfo struct {
	Gender           optional.String
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
	UpdateAt         time.Time
}
