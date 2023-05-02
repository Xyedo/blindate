package v1

import (
	"github.com/xyedo/blindate/internal/optional"
)

type getBasicInfoResponse struct {
	Gender           string          `json:"gender"`
	FromLoc          optional.String `json:"from_loc"`
	Height           optional.Int16  `json:"height"`
	EducationLevel   optional.String `json:"education_level"`
	Drinking         optional.String `json:"drinking"`
	Smoking          optional.String `json:"smoking"`
	RelationshipPref optional.String `json:"relationship_pref"`
	LookingFor       string          `json:"looking_for"`
	Zodiac           optional.String `json:"zodiac"`
	Kids             optional.Int16  `json:"kids"`
	Work             optional.String `json:"work"`
}
