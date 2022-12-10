package service

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xyedo/blindate/pkg/domain"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 20 * time.Second
	pingPeriod = 10 * time.Second
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
		ws.cleanUp(conn)
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))

	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	var payload domain.WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			log.Println(err)
			return
		} else {
			payload.Conn = *conn
			ws.WsChan <- payload
		}
		select {
		case <-pingTicker.C:
			err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait))
			if err != nil {
				log.Println("ping err", err)
				return
			}
		default:
			continue
		}

	}
}

func (ws *Ws) cleanUp(conn *domain.WsConn) {
	conn.Close()
	userId, ok := ws.Clients[*conn]
	if !ok {
		return
	}
	delete(ws.Clients, *conn)
	delete(ws.ReverseClient, userId)

}
