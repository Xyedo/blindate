package websocketEntity

import "github.com/gorilla/websocket"

type Response struct {
	Action string         `json:"action"`
	Data   map[string]any `json:"data"`
}

type Payload struct {
	Action  string `json:"action"`
	Payload string `json:"payload"`
	Conn    Conn   `json:"-"`
}
type Conn struct {
	*websocket.Conn
}
