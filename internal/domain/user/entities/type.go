package entities

type Gender string

const (
	GenderFemale Gender = "FEMALE"
	GenderMale   Gender = "MALE"
	GenderOther  Gender = "Other"
)

type EducationLevel string

const (
	EducationLevelBeforeHighSchool EducationLevel = "Less than high school diploma"
	EducationLevelHighSchool       EducationLevel = "High school"
	EducationLevelAttendCollege    EducationLevel = "Some college, no degree"
	EducationLevelAssociate        EducationLevel = "Assosiate''s Degree"
	EducationLevelBachelor         EducationLevel = "Bachelor''s Degree"
	EducationLevelMaster           EducationLevel = "Master''s Degree"
	EducationLevelProfessional     EducationLevel = "Professional Degree"
	EducationLevelDoctorate        EducationLevel = "Doctorate Degree"
)

type DrinkingLevel string

const (
	DrinkingLevelNever             DrinkingLevel = "Never"
	DrinkingLevelOccasionally      DrinkingLevel = "Ocassionally"
	DrinkingLevelOnceAWeek         DrinkingLevel = "Once a week"
	DrinkingLevelMoreThanOnceAWeek DrinkingLevel = "More than 2/3 times a week"
	DrinkingLevelEveryDay          DrinkingLevel = "Every day"
)

type SmokingLevel string

const (
	SmokingLevelNever             SmokingLevel = "Never"
	SmokingLevelOccasionally      SmokingLevel = "Ocassionally"
	SmokingLevelOnceAWeek         SmokingLevel = "Once a week"
	SmokingLevelMoreThanOnceAWeek SmokingLevel = "More than 2/3 times a week"
	SmokingLevelEveryDay          SmokingLevel = "Every day"
)

type RelationshipPreference string

const (
	RelationshipPreferenceONS     RelationshipPreference = "One night Stand"
	RelationshipPreferenceCasual  RelationshipPreference = "Casual"
	RelationshipPreferenceSerious RelationshipPreference = "Serious"
)

type Zodiac string

const (
	ZodiacAries       Zodiac = "Aries"
	ZodiacTaurus      Zodiac = "Taurus"
	ZodiacGemini      Zodiac = "Gemini"
	ZodiacCancer      Zodiac = "Cancer"
	ZodiacLeo         Zodiac = "Leo"
	ZodiacVirgo       Zodiac = "Virgo"
	ZodiacLibra       Zodiac = "Libra"
	ZodiacScorpio     Zodiac = "Scorpio"
	ZodiacSagittarius Zodiac = "Sagittarius"
	ZodiacCapricorn   Zodiac = "Capricorn"
	ZodiacAquarius    Zodiac = "Aquarius"
	ZodiacPisces      Zodiac = "Pisces"
)
