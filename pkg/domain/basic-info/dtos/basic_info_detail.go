package basicInfoDTOs

import (
	"time"

	"github.com/xyedo/blindate/internal/optional"
)

type BasicInfo struct {
	UserId           string          `db:"user_id"`
	Gender           string          `db:"gender"`
	FromLoc          optional.String `db:"from_loc"`
	Height           optional.Int16  `db:"height"`
	EducationLevel   optional.String `db:"education_level"`
	Drinking         optional.String `db:"drinking"`
	Smoking          optional.String `db:"smoking"`
	RelationshipPref optional.String `db:"relationship_pref"`
	LookingFor       string          `db:"looking_for"`
	Zodiac           optional.String `db:"zodiac"`
	Kids             optional.Int16  `db:"kids"`
	Work             optional.String `db:"work"`
	CreatedAt        time.Time       `db:"created_at"`
	UpdatedAt        time.Time       `db:"updated_at"`
}
