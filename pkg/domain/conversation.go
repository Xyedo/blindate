package domain

import "time"

type tinyUser struct {
	ID         string `db:"id" json:"id"`
	FullName   string `db:"full_name" json:"fullName"`
	Alias      string `db:"alias" json:"alias"`
	ProfilePic string `db:"pictureLink" json:"profilePicture,omitempty"`
}

//one user to many conv
type Conversation struct {
	Id                string    `json:"id,omitempty" db:"id"`
	FromUser          tinyUser  `json:"fromUser" db:"creator"`
	ToUser            tinyUser  `json:"toUser" db:"recipient"`
	LastMessage       string    `json:"lastMessageSent,omitempty" db:"last_message,omitempty"`
	LastMessageSentAt time.Time `json:"lastMessageSentAt,omitempty" db:"last_message_sent_at,omitempty"`
	LastMessageSeenAt time.Time `json:"lastMessageSeenAt,omitempty" db:"last_messsage_seen_at,omitempty"`
	ChatRows          int       `json:"chatRows" db:"chat_rows"`
	DayPass           int       `json:"dayPass" db:"day_pass"`
}

//one convoersation to many chat
type Chat struct {
	Id             string          `json:"id"`
	ConversationId string          `json:"conversationId"`
	Messages       string          `json:"messages"`
	ReplyTo        *string         `json:"replyTo"`
	SentAt         time.Time       `json:"sentAt"`
	SeenAt         *time.Time      `json:"seenAt"`
	Attachment     *ChatAttachment `json:"attachment"`
}

//one to one with chat
type ChatAttachment struct {
	ChatId    string `json:"-" db:"chat_id"`
	BlobLink  string `json:"blobLink" db:"blob_link"`
	MediaType string `json:"mediaType" db:"media_type"`
}
