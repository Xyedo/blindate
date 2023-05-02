package v1

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/xyedo/blindate/internal/optional"
	"github.com/xyedo/blindate/pkg/common/mod"
)

type postBasicInfoRequest struct {
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
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

func (b *postBasicInfoRequest) mod() *postBasicInfoRequest {
	if b == nil {
		return nil
	}

	mod.Trim(&b.Gender)

	b.FromLoc.If(func(fromLoc string) {
		mod.TrimWhiteSpace(&fromLoc)
		b.FromLoc.Set(fromLoc)
	})

	b.EducationLevel.If(func(educationLevel string) {
		mod.TrimWhiteSpace(&educationLevel)
		b.EducationLevel.Set(educationLevel)
	})

	b.Drinking.If(func(drinking string) {
		mod.TrimWhiteSpace(&drinking)
		b.Drinking.Set(drinking)
	})

	b.Smoking.If(func(smoking string) {
		mod.TrimWhiteSpace(&smoking)
		b.Smoking.Set(smoking)
	})

	b.RelationshipPref.If(func(relationshipPref string) {
		mod.TrimWhiteSpace(&relationshipPref)
		b.RelationshipPref.Set(relationshipPref)
	})

	mod.TrimWhiteSpace(&b.LookingFor)

	b.Zodiac.If(func(zodiac string) {
		mod.TrimWhiteSpace(&zodiac)
		b.Zodiac.Set(zodiac)
	})

	b.Work.If(func(work string) {
		mod.TrimWhiteSpace(&work)
		b.Work.Set(work)
	})
	return b
}

func (b postBasicInfoRequest) validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Gender, validation.Required, validation.In("Female", "Male", "Other")),
		validation.Field(&b.FromLoc, validation.Required.When(b.FromLoc.ValueSet()), validation.Length(1, 99)),
		validation.Field(&b.Height, validation.Required.When(b.Height.ValueSet()), validation.Min(0), validation.Max(99)),
		validation.Field(&b.EducationLevel, validation.Required.When(b.EducationLevel.ValueSet()), validation.In("Less than high school diploma", "High school", "Some college, no degree", "Assosiate's Degree", "Bachelor's Degree", "Master's Degree", "Professional Degree", "Doctorate Degree")),
		validation.Field(&b.Drinking, validation.Required.When(b.Drinking.ValueSet()), validation.In("Never Ocassionally", "Once a week", "More than 2/3 times a week", "Every day")),
		validation.Field(&b.Smoking, validation.Required.When(b.Smoking.ValueSet()), validation.In("Never Ocassionally", "Once a week", "More than 2/3 times a week", "Every day")),
		validation.Field(&b.RelationshipPref, validation.Required.When(b.RelationshipPref.ValueSet()), validation.In("One night Stand", "Having fun", "Serious")),
		validation.Field(&b.LookingFor, validation.Required, validation.In("Female", "Male", "Other")),
		validation.Field(&b.Zodiac, validation.Required.When(b.Zodiac.ValueSet()), validation.In("Aries", "Taurus", "Gemini", "Cancer", "Leo", "Virgo", "Libra", "Scorpio", "Sagittarius", "Capricorn", "Aquarius", "Pisces")),
		validation.Field(&b.Kids, validation.Required.When(b.Kids.ValueSet()), validation.Min(0), validation.Max(30)),
		validation.Field(&b.Work, validation.Required.When(b.Work.ValueSet()), validation.Length(5, 50)),
	)
}
