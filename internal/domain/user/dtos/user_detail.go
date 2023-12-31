package dtos

import (
	"github.com/invopop/validation"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/common/mod"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/pkg/optional"
)

type PostUserDetailRequest struct {
	Alias            string          `json:"alias"`
	Gender           string          `json:"gender"`
	Location         Location        `json:"location"`
	Bio              string          `json:"bio"`
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

func (req *PostUserDetailRequest) Mod() *PostUserDetailRequest {
	mod.Trim(&req.Bio)
	mod.TrimWhiteSpace(&req.Alias)

	return req
}

func (req PostUserDetailRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Alias, validation.Required, validation.Length(5, 200)),
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
		validation.Field(&req.LookingFor,
			validation.Required,
			validation.In(
				string(entities.GenderFemale),
				string(entities.GenderMale),
				string(entities.GenderOther),
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
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (l Location) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Lat, validation.Required, validation.Max(float64(90.0)), validation.Min(float64(-90.0))),
		validation.Field(&l.Lng, validation.Required, validation.Min(float64(-180.0)), validation.Max(float64(180.0))),
	)
}

type PatchUserDetailRequest struct {
	Alias            optional.String           `json:"alias"`
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
	if !req.Alias.IsSet() &&
		!req.Gender.IsSet() &&
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
		validation.Field(&req.Alias, validation.Required.When(req.Alias.IsSet()), validation.Length(5, 200)),
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
		validation.Field(&req.LookingFor,
			validation.Required.When(req.Alias.IsSet()),
			validation.In(
				string(entities.GenderFemale),
				string(entities.GenderMale),
				string(entities.GenderOther),
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

type UserInterest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserDetail struct {
	UserId           string              `json:"user_id"`
	Alias            string              `json:"alias"`
	Geog             UserDetailGeography `json:"geo"`
	Bio              string              `json:"bio"`
	Gender           string              `json:"gender"`
	FromLoc          optional.String     `json:"from_location"`
	Height           optional.Int16      `json:"height"`
	EducationLevel   optional.String     `json:"education_level"`
	Drinking         optional.String     `json:"drinking"`
	Smoking          optional.String     `json:"smoking"`
	RelationshipPref optional.String     `json:"relationship_preferences"`
	LookingFor       string              `json:"looking_for"`
	Zodiac           optional.String     `json:"zodiac"`
	Kids             optional.Int16      `json:"kids"`
	Work             optional.String     `json:"work"`

	Hobbies            []UserInterest `json:"hobbies"`
	MovieSeries        []UserInterest `json:"movie_series"`
	Travels            []UserInterest `json:"travels"`
	Sports             []UserInterest `json:"sports"`
	ProfilePictureURLs []string       `json:"profile_picture_urls"`
}
type UserDetailGeography struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func NewUserDetailResponse(user entities.UserDetail) UserDetail {
	hobbies := make([]UserInterest, 0, len(user.Hobbies))
	for _, hobbie := range user.Hobbies {
		hobbies = append(hobbies, UserInterest{
			Id:   hobbie.Id,
			Name: hobbie.Hobbie,
		})
	}

	movieSeries := make([]UserInterest, 0, len(user.MovieSeries))
	for _, movieSerie := range user.MovieSeries {
		movieSeries = append(movieSeries, UserInterest{
			Id:   movieSerie.Id,
			Name: movieSerie.MovieSerie,
		})
	}

	travels := make([]UserInterest, 0, len(user.Travels))
	for _, travel := range user.Travels {
		travels = append(travels, UserInterest{
			Id:   travel.Id,
			Name: travel.Travel,
		})
	}

	sports := make([]UserInterest, 0, len(user.Sports))
	for _, sport := range user.Sports {
		sports = append(sports, UserInterest{
			Id:   sport.Id,
			Name: sport.Sport,
		})
	}

	profilePicURLs := make([]string, 0, len(user.ProfilePictures))
	for _, profilePic := range user.ProfilePictures {
		profilePicURLs = append(profilePicURLs, profilePic.GetPresignedUrl())
	}

	return UserDetail{
		UserId: user.UserId,
		Alias:  user.Alias,
		Geog: UserDetailGeography{
			Lat: user.Geog.Lat,
			Lng: user.Geog.Lng,
		},
		Bio:                user.Bio,
		Gender:             string(user.Gender),
		FromLoc:            user.FromLoc,
		Height:             user.Height,
		EducationLevel:     user.EducationLevel,
		Drinking:           user.Drinking,
		Smoking:            user.Smoking,
		RelationshipPref:   user.RelationshipPref,
		LookingFor:         user.LookingFor,
		Zodiac:             user.Zodiac,
		Kids:               user.Kids,
		Work:               user.Work,
		Hobbies:            hobbies,
		MovieSeries:        movieSeries,
		Travels:            travels,
		Sports:             sports,
		ProfilePictureURLs: profilePicURLs,
	}
}
