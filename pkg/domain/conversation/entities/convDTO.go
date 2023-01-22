package convEntity

import "time"

type tinyUser struct {
	ID         string `db:"id" json:"id"`
	FullName   string `db:"full_name" json:"fullName,omitempty"`
	Alias      string `db:"alias" json:"alias,omitempty"`
	ProfilePic string `db:"picture_ref" json:"profilePicture,omitempty"`
}

// Conversation one to one with match
type DTO struct {
	Id                string     `json:"id,omitempty" db:"id"`
	FromUser          tinyUser   `json:"fromUser" db:"creator"`
	ToUser            tinyUser   `json:"toUser" db:"recipient"`
	LastMessage       string     `json:"lastMessageSent,omitempty" db:"last_message"`
	LastMessageSentAt time.Time  `json:"lastMessageSentAt,omitempty" db:"last_message_sent_at"`
	LastMessageSeenAt *time.Time `json:"lastMessageSeenAt,omitempty" db:"last_messsage_seen_at"`
	ChatRows          int        `json:"chatRows" db:"chat_rows"`
	DayPass           int        `json:"dayPass" db:"day_pass"`
	RequestStatus     string     `json:"-"`
	RevealStatus      string     `json:"-"`
}
