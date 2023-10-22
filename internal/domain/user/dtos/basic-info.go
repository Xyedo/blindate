package dtos

import (
	"github.com/invopop/validation"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/pkg/optional"
)

type PostBasicInfoRequest struct {
	Gender           string          `json:"gender"`
	FromLoc          optional.String `json:"from_loc"`
	Height           optional.Int16  `json:"height"`
	EducationLevel   optional.String `json:"education_level"`
	Drinking         optional.String `json:"drinking"`
	Smoking          optional.String `json:"smoking"`
	RelationshipPref optional.String `json:"relationship_pref"`
	LookingFor       optional.String `json:"looking_for"`
	Zodiac           optional.String `json:"zodiac"`
	Kids             optional.Int16  `json:"kids"`
	Work             optional.String `json:"work"`
}

func (req PostBasicInfoRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Gender, validation.Required,
			validation.In(
				entities.GenderFemale,
				entities.GenderMale,
				entities.GenderOther,
			),
		),
		validation.Field(&req.FromLoc, validation.Length(0, 100)),
		validation.Field(&req.Height, validation.Max(400)),
		validation.Field(&req.EducationLevel,
			validation.In(
				entities.EducationLevelBeforeHighSchool,
				entities.EducationLevelHighSchool,
				entities.EducationLevelAttendCollege,
				entities.EducationLevelAssociate,
				entities.EducationLevelBachelor,
				entities.EducationLevelMaster,
				entities.EducationLevelProfessional,
				entities.EducationLevelDoctorate,
			),
		),
		validation.Field(&req.Drinking,
			validation.In(
				entities.DrinkingLevelNever,
				entities.DrinkingLevelOccasionally,
				entities.DrinkingLevelOnceAWeek,
				entities.DrinkingLevelMoreThanOnceAWeek,
				entities.DrinkingLevelEveryDay,
			),
		),
		validation.Field(&req.Smoking,
			validation.In(
				entities.SmokingLevelNever,
				entities.SmokingLevelOccasionally,
				entities.SmokingLevelOnceAWeek,
				entities.SmokingLevelMoreThanOnceAWeek,
				entities.SmokingLevelEveryDay,
			),
		),
		validation.Field(&req.RelationshipPref,
			validation.In(
				entities.RelationshipPreferenceONS,
				entities.RelationshipPreferenceCasual,
				entities.RelationshipPreferenceSerious,
			),
		),
		validation.Field(&req.Zodiac,
			validation.In(
				entities.ZodiacAries,
				entities.ZodiacTaurus,
				entities.ZodiacGemini,
				entities.ZodiacCancer,
				entities.ZodiacLeo,
				entities.ZodiacVirgo,
				entities.ZodiacLibra,
				entities.ZodiacScorpio,
				entities.ZodiacSagittarius,
				entities.ZodiacCapricorn,
				entities.ZodiacAquarius,
				entities.ZodiacPisces,
			),
		),
		validation.Field(&req.Kids, validation.Max(100)),
		validation.Field(&req.Work, validation.Length(0, 50)),
	)
}

type PatchBasicInfoRequest struct {
	Gender           optional.String `json:"gender"`
	FromLoc          optional.String `json:"from_loc"`
	Height           optional.Int16  `json:"height"`
	EducationLevel   optional.String `json:"education_level"`
	Drinking         optional.String `json:"drinking"`
	Smoking          optional.String `json:"smoking"`
	RelationshipPref optional.String `json:"relationship_pref"`
	LookingFor       optional.String `json:"looking_for"`
	Zodiac           optional.String `json:"zodiac"`
	Kids             optional.Int16  `json:"kids"`
	Work             optional.String `json:"work"`
}

func (req PatchBasicInfoRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Gender, validation.Required.When(req.Gender.IsSet()),
			validation.In(
				entities.GenderFemale,
				entities.GenderMale,
				entities.GenderOther,
			),
		),
		validation.Field(&req.FromLoc, validation.Length(0, 100)),
		validation.Field(&req.Height, validation.Max(400)),
		validation.Field(&req.EducationLevel,
			validation.In(
				entities.EducationLevelBeforeHighSchool,
				entities.EducationLevelHighSchool,
				entities.EducationLevelAttendCollege,
				entities.EducationLevelAssociate,
				entities.EducationLevelBachelor,
				entities.EducationLevelMaster,
				entities.EducationLevelProfessional,
				entities.EducationLevelDoctorate,
			),
		),
		validation.Field(&req.Drinking,
			validation.In(
				entities.DrinkingLevelNever,
				entities.DrinkingLevelOccasionally,
				entities.DrinkingLevelOnceAWeek,
				entities.DrinkingLevelMoreThanOnceAWeek,
				entities.DrinkingLevelEveryDay,
			),
		),
		validation.Field(&req.Smoking,
			validation.In(
				entities.SmokingLevelNever,
				entities.SmokingLevelOccasionally,
				entities.SmokingLevelOnceAWeek,
				entities.SmokingLevelMoreThanOnceAWeek,
				entities.SmokingLevelEveryDay,
			),
		),
		validation.Field(&req.RelationshipPref,
			validation.In(
				entities.RelationshipPreferenceONS,
				entities.RelationshipPreferenceCasual,
				entities.RelationshipPreferenceSerious,
			),
		),
		validation.Field(&req.Zodiac,
			validation.In(
				entities.ZodiacAries,
				entities.ZodiacTaurus,
				entities.ZodiacGemini,
				entities.ZodiacCancer,
				entities.ZodiacLeo,
				entities.ZodiacVirgo,
				entities.ZodiacLibra,
				entities.ZodiacScorpio,
				entities.ZodiacSagittarius,
				entities.ZodiacCapricorn,
				entities.ZodiacAquarius,
				entities.ZodiacPisces,
			),
		),
		validation.Field(&req.Kids, validation.Max(100)),
		validation.Field(&req.Work, validation.Length(0, 50)),
	)
}
