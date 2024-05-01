package handler

import "github.com/xyedo/blindate/internal/domain/conversation"

func New(conversationUsecase conversation.Usecase) *Conversation {
	return &Conversation{
		usecase: conversationUsecase,
	}
}

type Conversation struct {
	usecase conversation.Usecase
}
