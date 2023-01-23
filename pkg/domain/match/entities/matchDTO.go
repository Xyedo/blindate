package matchEntity

import "time"

// Match one user to many match
type MatchDTO struct {
	Id            string     `json:"id"`
	RequestFrom   string     `json:"requestFrom"`
	RequestTo     string     `json:"requestTo"`
	RequestStatus Status     `json:"requestStatus"`
	CreatedAt     time.Time  `json:"createdAt"`
	AcceptedAt    *time.Time `json:"acceptedAt,omitempty"`
	RevealStatus  Status     `json:"revealStatus,omitempty"`
	RevealedAt    *time.Time `json:"revealedAt,omitempty"`
}
