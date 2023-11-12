package dtos

import (
	"strings"

	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/pkg/optional"
)

type PostUserDetailRequest struct {
	Gender           string          `json:"gender"`
	Location         Location        `json:"location"`
	Bio              string          `json:"bio"`
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

func (req *PostUserDetailRequest) Mod() *PostUserDetailRequest {
	req.Bio = strings.TrimSpace(req.Bio)

	return req
}

func (req PostUserDetailRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Gender, validation.Required,
			validation.In(
				string(entities.GenderFemale),
				string(entities.GenderMale),
				string(entities.GenderOther),
			),
		),
		validation.Field(&req.Location),
		validation.Field(&req.Bio, validation.Required, validation.Length(2, 300)),
		validation.Field(&req.FromLoc, validation.Length(0, 100)),
		validation.Field(&req.Height, validation.Max(400)),
		validation.Field(&req.EducationLevel,
			validation.In(
				string(entities.EducationLevelBeforeHighSchool),
				string(entities.EducationLevelHighSchool),
				string(entities.EducationLevelAttendCollege),
				string(entities.EducationLevelAssociate),
				string(entities.EducationLevelBachelor),
				string(entities.EducationLevelMaster),
				string(entities.EducationLevelProfessional),
				string(entities.EducationLevelDoctorate),
			),
		),
		validation.Field(&req.Drinking,
			validation.In(
				string(entities.DrinkingLevelNever),
				string(entities.DrinkingLevelOccasionally),
				string(entities.DrinkingLevelOnceAWeek),
				string(entities.DrinkingLevelMoreThanOnceAWeek),
				string(entities.DrinkingLevelEveryDay),
			),
		),
		validation.Field(&req.Smoking,
			validation.In(
				string(entities.SmokingLevelNever),
				string(entities.SmokingLevelOccasionally),
				string(entities.SmokingLevelOnceAWeek),
				string(entities.SmokingLevelMoreThanOnceAWeek),
				string(entities.SmokingLevelEveryDay),
			),
		),
		validation.Field(&req.RelationshipPref,
			validation.In(
				string(entities.RelationshipPreferenceONS),
				string(entities.RelationshipPreferenceCasual),
				string(entities.RelationshipPreferenceSerious),
			),
		),
		validation.Field(&req.Zodiac,
			validation.In(
				string(entities.ZodiacAries),
				string(entities.ZodiacTaurus),
				string(entities.ZodiacGemini),
				string(entities.ZodiacCancer),
				string(entities.ZodiacLeo),
				string(entities.ZodiacVirgo),
				string(entities.ZodiacLibra),
				string(entities.ZodiacScorpio),
				string(entities.ZodiacSagittarius),
				string(entities.ZodiacCapricorn),
				string(entities.ZodiacAquarius),
				string(entities.ZodiacPisces),
			),
		),
		validation.Field(&req.Kids, validation.Max(100)),
		validation.Field(&req.Work, validation.Length(0, 50)),
	)
}

type Location struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

func (l Location) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Lat, validation.Required, is.Latitude),
		validation.Field(&l.Lng, validation.Required, is.Longitude),
	)
}

type PatchUserDetailRequest struct {
	Gender           optional.String           `json:"gender"`
	Location         optional.Option[Location] `json:"location"`
	Bio              optional.String           `json:"bio"`
	FromLoc          optional.String           `json:"from_loc"`
	Height           optional.Int16            `json:"height"`
	EducationLevel   optional.String           `json:"education_level"`
	Drinking         optional.String           `json:"drinking"`
	Smoking          optional.String           `json:"smoking"`
	RelationshipPref optional.String           `json:"relationship_pref"`
	LookingFor       optional.String           `json:"looking_for"`
	Zodiac           optional.String           `json:"zodiac"`
	Kids             optional.Int16            `json:"kids"`
	Work             optional.String           `json:"work"`
}

func (req PatchUserDetailRequest) Validate() error {
	if !req.Gender.IsSet() &&
		!req.Location.IsSet() &&
		!req.Bio.IsSet() &&
		!req.FromLoc.IsSet() &&
		!req.Height.IsSet() &&
		!req.EducationLevel.IsSet() &&
		!req.Drinking.IsSet() &&
		!req.Smoking.IsSet() &&
		!req.RelationshipPref.IsSet() &&
		!req.LookingFor.IsSet() &&
		!req.Zodiac.IsSet() &&
		!req.Kids.IsSet() &&
		!req.Work.IsSet() {
		return apperror.BadPayload(apperror.Payload{
			Status:  apperror.StatusErrorValidation,
			Message: "body shouldn't empty",
		})
	}
	return validation.ValidateStruct(&req,
		validation.Field(&req.Gender, validation.Required.When(req.Gender.IsSet()),
			validation.In(
				string(entities.GenderFemale),
				string(entities.GenderMale),
				string(entities.GenderOther),
			),
		),
		validation.Field(&req.Location, validation.Skip.When(!req.Location.IsSet())),
		validation.Field(&req.Bio, validation.Required.When(req.Bio.IsSet()), validation.Length(2, 300)),
		validation.Field(&req.FromLoc, validation.Length(0, 100)),
		validation.Field(&req.Height, validation.Max(400)),
		validation.Field(&req.EducationLevel,
			validation.In(
				string(entities.EducationLevelBeforeHighSchool),
				string(entities.EducationLevelHighSchool),
				string(entities.EducationLevelAttendCollege),
				string(entities.EducationLevelAssociate),
				string(entities.EducationLevelBachelor),
				string(entities.EducationLevelMaster),
				string(entities.EducationLevelProfessional),
				string(entities.EducationLevelDoctorate),
			),
		),
		validation.Field(&req.Drinking,
			validation.In(
				string(entities.DrinkingLevelNever),
				string(entities.DrinkingLevelOccasionally),
				string(entities.DrinkingLevelOnceAWeek),
				string(entities.DrinkingLevelMoreThanOnceAWeek),
				string(entities.DrinkingLevelEveryDay),
			),
		),
		validation.Field(&req.Smoking,
			validation.In(
				string(entities.SmokingLevelNever),
				string(entities.SmokingLevelOccasionally),
				string(entities.SmokingLevelOnceAWeek),
				string(entities.SmokingLevelMoreThanOnceAWeek),
				string(entities.SmokingLevelEveryDay),
			),
		),
		validation.Field(&req.RelationshipPref,
			validation.In(
				string(entities.RelationshipPreferenceONS),
				string(entities.RelationshipPreferenceCasual),
				string(entities.RelationshipPreferenceSerious),
			),
		),
		validation.Field(&req.Zodiac,
			validation.In(
				string(entities.ZodiacAries),
				string(entities.ZodiacTaurus),
				string(entities.ZodiacGemini),
				string(entities.ZodiacCancer),
				string(entities.ZodiacLeo),
				string(entities.ZodiacVirgo),
				string(entities.ZodiacLibra),
				string(entities.ZodiacScorpio),
				string(entities.ZodiacSagittarius),
				string(entities.ZodiacCapricorn),
				string(entities.ZodiacAquarius),
				string(entities.ZodiacPisces),
			),
		),
		validation.Field(&req.Kids, validation.Max(100)),
		validation.Field(&req.Work, validation.Length(0, 50)),
	)
}

func (req PatchUserDetailRequest) ToEntity() entities.UpdateUserDetail {
	var geog optional.Option[entities.Geography]
	req.Location.If(func(l Location) {
		geog = optional.New(entities.Geography{
			Lat: l.Lat,
			Lng: l.Lng,
		})
	})

	return entities.UpdateUserDetail{
		Gender:           req.Gender,
		Geog:             geog,
		Bio:              req.Bio,
		FromLoc:          req.FromLoc,
		Height:           req.Height,
		EducationLevel:   req.EducationLevel,
		Drinking:         req.Drinking,
		Smoking:          req.Smoking,
		RelationshipPref: req.RelationshipPref,
		LookingFor:       req.LookingFor,
		Zodiac:           req.Zodiac,
		Kids:             req.Kids,
		Work:             req.Work,
	}
}

type UserDetail struct {
	UserId           string              `json:"user_id"`
	Geog             UserDetailGeography `json:"geo"`
	Bio              string              `json:"bio"`
	Gender           string              `json:"gender"`
	FromLoc          optional.String     `json:"from_location"`
	Height           optional.Int16      `json:"height"`
	EducationLevel   optional.String     `json:"education_level"`
	Drinking         optional.String     `json:"drinking"`
	Smoking          optional.String     `json:"smoking"`
	RelationshipPref optional.String     `json:"relationship_preferences"`
	LookingFor       optional.String     `json:"looking_for"`
	Zodiac           optional.String     `json:"zodiac"`
	Kids             optional.Int16      `json:"kids"`
	Work             optional.String     `json:"work"`

	Hobbies     []string `json:"hobbies"`
	MovieSeries []string `json:"movie_series"`
	Travels     []string `json:"travels"`
	Sports      []string `json:"sports"`
}
type UserDetailGeography struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

func GetUserDetailResponse(user entities.UserDetail) UserDetail {
	hobbies := make([]string, 0, len(user.Hobbies))
	for _, hobbie := range user.Hobbies {
		hobbies = append(hobbies, hobbie.Hobbie)
	}

	movieSeries := make([]string, 0, len(user.MovieSeries))
	for _, movieSerie := range user.MovieSeries {
		movieSeries = append(movieSeries, movieSerie.MovieSerie)
	}

	travels := make([]string, 0, len(user.Travels))
	for _, travel := range user.Travels {
		travels = append(travels, travel.Travel)
	}

	sports := make([]string, 0, len(user.Sports))
	for _, sport := range user.Sports {
		sports = append(sports, sport.Sport)
	}

	return UserDetail{
		UserId: user.UserId,
		Geog: UserDetailGeography{
			Lat: user.Geog.Lat,
			Lng: user.Geog.Lng,
		},
		Bio:              user.Bio,
		Gender:           string(user.Gender),
		FromLoc:          user.FromLoc,
		Height:           user.Height,
		EducationLevel:   user.EducationLevel,
		Drinking:         user.Drinking,
		Smoking:          user.Smoking,
		RelationshipPref: user.RelationshipPref,
		LookingFor:       user.LookingFor,
		Zodiac:           user.Zodiac,
		Kids:             user.Kids,
		Work:             user.Work,
		Hobbies:          hobbies,
		MovieSeries:      movieSeries,
		Travels:          travels,
		Sports:           sports,
	}
}
