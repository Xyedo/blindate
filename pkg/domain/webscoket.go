package domain

import "github.com/gorilla/websocket"

type WsResponse struct {
	Action string         `json:"action"`
	Data   map[string]any `json:"data"`
}

type WsPayload struct {
	Action  string `json:"action"`
	Payload string `json:"payload"`
	Conn    WsConn `json:"-"`
}
type WsConn struct {
	*websocket.Conn
}
