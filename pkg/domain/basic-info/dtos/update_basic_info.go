package basicInfoDTOs

import (
	"time"

	"github.com/xyedo/blindate/internal/optional"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
)

type UpdateBasicInfo struct {
	UserId           string
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
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (b UpdateBasicInfo) Validate() error {
	if !b.Gender.ValueSet() &&
		!b.FromLoc.ValueSet() &&
		!b.Height.ValueSet() &&
		!b.EducationLevel.ValueSet() &&
		!b.Drinking.ValueSet() &&
		!b.Smoking.ValueSet() &&
		!b.RelationshipPref.ValueSet() &&
		!b.LookingFor.ValueSet() &&
		!b.Zodiac.ValueSet() &&
		!b.Kids.ValueSet() &&
		!b.Work.ValueSet() {
		return apperror.BadPayload(apperror.Payload{Message: "body should not be empty"})
	}
	return nil
}
