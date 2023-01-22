package basicInfoEntity

import "time"

// BasicInfo one to one with user
type FullDTO struct {
	UserId           string    `json:"userId"`
	Gender           string    `json:"gender" binding:"required,oneof=Female Male Other"`
	FromLoc          *string   `json:"fromLoc" binding:"omitempty,max=99"`
	Height           *int      `json:"height" binding:"omitempty,min=0,max=400"`
	EducationLevel   *string   `json:"educationLevel" binding:"omitempty,valideducationlevel"`
	Drinking         *string   `json:"drinking" binding:"omitempty,oneof=Never Ocassionally 'Once a week' 'More than 2/3 times a week' 'Every day'"`
	Smoking          *string   `json:"smoking" binding:"omitempty,oneof=Never Ocassionally 'Once a week' 'More than 2/3 times a week' 'Every day'"`
	RelationshipPref *string   `json:"relationshipPref" binding:"omitempty,oneof='One night Stand' 'Having fun' Serious"`
	LookingFor       string    `json:"lookingFor" binding:"required,oneof=Female Male Other"`
	Zodiac           *string   `json:"zodiac" binding:"omitempty,oneof=Aries Taurus Gemini Cancer Leo Virgo Libra Scorpio Sagittarius Capricorn Aquarius Pisces"`
	Kids             *int      `json:"kids" binding:"omitempty,min=0,max=30"`
	Work             *string   `json:"work" binding:"omitempty,max=50"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
