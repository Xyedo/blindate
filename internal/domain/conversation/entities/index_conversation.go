package entities

import (
	"time"

	"github.com/xyedo/blindate/pkg/optional"
)

type ConversationIndex []ConversationElement

type ConversationElement struct {
	Id        string               `json:"id"`
	Recepient RecepientUser        `json:"recepient"`
	ChatRows  int64                `json:"chat_rows"`
	DayPass   int64                `json:"day_pass"`
	LastChat  ConversationLastChat `json:"last_chat"`
	UpdatedAt time.Time            `json:"updated_at"`
	CreatedAt time.Time            `json:"created_at"`
	Version   int64                `json:"-"`
}

type RecepientUser struct {
	Id          string          `json:"id"`
	DisplayName string          `json:"display_name"`
	FileId      optional.String `json:"-"`
	Url         string          `json:"url"`
}
type ConversationLastChat struct {
	Id                 optional.String `json:"id"`
	Author             optional.String `json:"author"`
	Message            optional.String `json:"message"`
	UnreadMessageCount optional.Int64  `json:"unread_message_count"`
	ReplyTo            optional.String `json:"reply_to"`
	SentAt             optional.Time   `json:"sent_at"`
	SeenAt             optional.Time   `json:"seen_at"`
	UpdatedAt          optional.Time   `json:"updated_at"`
	Version            optional.Int64  `json:"-"`
}

func (convos ConversationIndex) ToFileIds() ([]string, map[string]int) {
	fileIds := make([]string, 0, len(convos))
	fieIdToIdx := make(map[string]int, 0)

	for i, convo := range convos {
		convo.Recepient.FileId.If(func(fileId string) {
			fileIds = append(fileIds, fileId)
			fieIdToIdx[fileId] = i
		})
	}

	return fileIds, fieIdToIdx
}
