package event

type ConnectionPayload struct {
	UserId string
	Online bool
}

var UserConnectionChange userConnectionChange

type userConnectionChange struct {
	handlers []interface{ HandleUserChangeConnection(ConnectionPayload) }
}

func (u *userConnectionChange) Register(handler interface{ HandleUserChangeConnection(ConnectionPayload) }) {
	u.handlers = append(u.handlers, handler)
}

func (u userConnectionChange) Trigger(payload ConnectionPayload) {
	for _, handler := range u.handlers {
		go handler.HandleUserChangeConnection(payload)
	}
}

func (u userConnectionChange) TriggerSync(payload ConnectionPayload) {
	for _, handler := range u.handlers {
		handler.HandleUserChangeConnection(payload)
	}
}
