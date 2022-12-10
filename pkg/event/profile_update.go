package event

var ProfileUpdated profileUpdated

type ProfileUpdatedPayload struct {
	UserId string
}
type profileUpdated struct {
	handlers []interface{ HandleProfileUpdateEvent(ProfileUpdatedPayload) }
}

func (u *profileUpdated) Register(handler interface{ HandleProfileUpdateEvent(ProfileUpdatedPayload) }) {
	u.handlers = append(u.handlers, handler)
}
func (u profileUpdated) Trigger(payload ProfileUpdatedPayload) {
	for _, handler := range u.handlers {
		go handler.HandleProfileUpdateEvent(payload)
	}
}
