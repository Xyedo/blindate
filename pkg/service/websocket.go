package service

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xyedo/blindate/internal/rwmap"
	"github.com/xyedo/blindate/pkg/domain"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 20 * time.Second
	pingPeriod = 10 * time.Second
)

func NewWs() *Ws {
	return &Ws{
		Clients:       rwmap.New[domain.WsConn, string](),
		ReverseClient: rwmap.New[string, domain.WsConn](),
		WsChan:        make(chan domain.WsPayload),
	}
}

type Ws struct {
	Clients       *rwmap.RwMap[domain.WsConn, string]
	ReverseClient *rwmap.RwMap[string, domain.WsConn]
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

	userId, ok := ws.Clients.Get(*conn)
	if !ok {
		return
	}
	ws.Clients.Delete(*conn)
	ws.ReverseClient.Delete(userId)

}
