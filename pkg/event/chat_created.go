package event

import "github.com/xyedo/blindate/pkg/domain"

var ChatCreated chatCreated

type ChatCreatedPayload struct {
	Chat   []domain.Chat
	ConvId string
}

type chatCreated struct {
	handlers []interface{ HandleCreateChatEvent(ChatCreatedPayload) }
}

func (m *chatCreated) Register(handler interface{ HandleCreateChatEvent(ChatCreatedPayload) }) {
	m.handlers = append(m.handlers, handler)
}

func (m chatCreated) Trigger(payload ChatCreatedPayload) {
	for _, handler := range m.handlers {
		go handler.HandleCreateChatEvent(payload)
	}
}
