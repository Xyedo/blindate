package event

var ChatSeen chatSeen

type ChatSeenPayload struct {
	RequestFrom string
	RequestTo   string
	SeenChatIds []string
}

type chatSeen struct {
	handlers []interface{ HandleSeenAtevent(ChatSeenPayload) }
}

func (m *chatSeen) Register(handler interface{ HandleSeenAtevent(ChatSeenPayload) }) {
	m.handlers = append(m.handlers, handler)
}

func (m chatSeen) Trigger(payload ChatSeenPayload) {
	for _, handler := range m.handlers {
		go handler.HandleSeenAtevent(payload)
	}
}
