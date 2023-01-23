package matchEntity

import (
	"time"

	interestEntity "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

type UserDTO struct {
	UserId           string                         ` json:"userId"`
	Alias            string                         ` json:"alias"`
	Dob              time.Time                      ` json:"dob"`
	Gender           *string                        `json:"gender"`
	FromLoc          *string                        `json:"fromLoc"`
	Height           *int                           `json:"height"`
	EducationLevel   *string                        `json:"educationLevel"`
	Drinking         *string                        `json:"drinking"`
	Smoking          *string                        `json:"smoking"`
	RelationshipPref *string                        `json:"relationshipPref"`
	LookingFor       *string                        `json:"lookingFor"`
	Zodiac           *string                        `json:"zodiac"`
	Kids             *int                           `json:"kids"`
	Work             *string                        `json:"work"`
	BioId            *string                        `json:"bioId"`
	Bio              *string                        `json:"bio" db:"bio"`
	Hobbies          []interestEntity.HobbieDTO     `json:"hobbies"`
	MovieSeries      []interestEntity.MovieSerieDTO `json:"movieSeries"`
	Travels          []interestEntity.TravelDTO     `json:"travels"`
	Sports           []interestEntity.SportDTO      `json:"sports"`
}
