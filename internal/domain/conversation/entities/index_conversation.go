package entities

import (
	"time"

	"github.com/xyedo/blindate/pkg/optional"
)

type ConversationIndex []ConversationElement

type ConversationElement struct {
	Id        string
	Recepient RecepientUser
	ChatRows  int64
	DayPass   int64
	LastChat  ConversationLastChat
	UpdatedAt time.Time
	CreatedAt time.Time
	Version   int64
}

type RecepientUser struct {
	Id          string
	DisplayName string
	FileId      optional.String
	Url         string
}
type ConversationLastChat struct {
	Id                 optional.String
	Author             optional.String
	Message            optional.String
	UnreadMessageCount optional.Int64
	ReplyTo            optional.String
	SentAt             optional.Time
	SeenAt             optional.Time
	UpdatedAt          optional.Time
	Version            optional.Int64
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
