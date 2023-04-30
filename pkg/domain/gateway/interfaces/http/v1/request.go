package v1

type requestPayload struct {
	Action  string `json:"action"`
	Payload string `json:"data"`
	UserId  string `json:"-"`
}
