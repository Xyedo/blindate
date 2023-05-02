package basicInfoDTOs

import (
	"time"

	"github.com/xyedo/blindate/internal/optional"
)

type CreateBasicInfo struct {
	UserId           string
	Gender           string
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
}
