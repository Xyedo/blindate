package v1

type response struct {
	Action string         `json:"action"`
	Data   map[string]any `json:"data"`
}
