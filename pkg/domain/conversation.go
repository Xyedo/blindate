package domain

import "time"

type Conversation struct {
	Id                string    `json:"id,omitempty" db:"id"`
	FromUser          User      `json:"fromUser" db:"creator"`
	ToUser            User      `json:"toUser" db:"recipient"`
	LastMessage       string    `json:"lastMessageSent,omitempty" db:"last_message,omitempty"`
	LastMessageSentAt time.Time `json:"lastMessageSentAt,omitempty" db:"last_message_sent_at,omitempty"`
	LastMessageSeenAt time.Time `json:"lastMessageSeenAt,omitempty" db:"last_messsage_seen_at,omitempty"`
	ChatRows          int       `json:"chatRows" db:"chat_rows"`
	DayPass           int       `json:"dayPass" db:"day_pass"`
}
