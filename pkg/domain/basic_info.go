package domain

import "time"

type BasicInfo struct {
	UserId           string    `json:"userId"`
	Gender           string    `json:"gender"`
	FromLoc          *string   `json:"fromLoc"`
	Height           *int      `json:"height"`
	EducationLevel   *string   `json:"educationLevel"`
	Drinking         *string   `json:"drinking"`
	Smoking          *string   `json:"smoking"`
	RelationshipPref *string   `json:"relationshipPref"`
	LookingFor       string    `json:"lookingFor"`
	Zodiac           *string   `json:"zodiac"`
	Kids             *int      `json:"kids"`
	Work             *string   `json:"work"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
