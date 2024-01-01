package entities

type Conversations []ConversationElement

type ConversationElement struct {
	Id             string
	DisplayName    string
	ProfilePic     string
	LastMessage    string
	UnreadMessages int
}
