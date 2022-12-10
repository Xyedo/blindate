package service

import (
	"log"

	"github.com/xyedo/blindate/pkg/domain"
)

func NewWs() *Ws {
	return &Ws{
		Clients:       make(map[domain.WsConn]string),
		ReverseClient: make(map[string]domain.WsConn),
		WsChan:        make(chan domain.WsPayload),
	}
}

type Ws struct {
	Clients       map[domain.WsConn]string
	ReverseClient map[string]domain.WsConn
	WsChan        chan domain.WsPayload
}

func (ws *Ws) ListenForWsPayload(conn *domain.WsConn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error", err)
		}
	}()
	var payload domain.WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {

		} else {
			payload.Conn = *conn
			ws.WsChan <- payload
		}
	}
}
