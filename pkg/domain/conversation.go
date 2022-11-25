package domain

import (
	"time"
)

type tinyUser struct {
	ID         string `db:"id" json:"id"`
	FullName   string `db:"full_name" json:"fullName,omitempty"`
	Alias      string `db:"alias" json:"alias,omitempty"`
	ProfilePic string `db:"picture_ref" json:"profilePicture,omitempty"`
}

// one user to many match
type Match struct {
	Id            string      `json:"id"`
	RequestFrom   string      `json:"requestFrom"`
	RequestTo     string      `json:"requestTo"`
	RequestStatus MatchStatus `json:"requestStatus"`
	CreatedAt     time.Time   `json:"createdAt"`
	AcceptedAt    *time.Time  `json:"acceptedAt,omitempty"`
	RevealStatus  MatchStatus `json:"revealStatus,omitempty"`
	RevealedAt    *time.Time  `json:"revealedAt,omitempty"`
}

// one to one with match
type Conversation struct {
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

// one convoersation to many chat
type Chat struct {
	Id             string          `json:"id"`
	ConversationId string          `json:"conversationId"`
	Author         string          `json:"author"`
	Messages       string          `json:"messages"`
	ReplyTo        *string         `json:"replyTo"`
	SentAt         time.Time       `json:"sentAt"`
	SeenAt         *time.Time      `json:"seenAt"`
	Attachment     *ChatAttachment `json:"attachment"`
}

// one to one with chat
type ChatAttachment struct {
	ChatId    string `json:"-" db:"chat_id"`
	BlobLink  string `json:"blobLink" db:"blob_link"`
	MediaType string `json:"mediaType" db:"media_type"`
}
