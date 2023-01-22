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
			log.Println("Error panic", err)
		}
		ws.cleanUp(conn)
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))

	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	var payload domain.WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}
		payload.Conn = *conn
		ws.WsChan <- payload
	}
}

func (ws *Ws) PingTicker(conn *domain.WsConn) {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		conn.Close()
	}()
	for {
		<-pingTicker.C
		err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeWait))
		if err != nil {
			return
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
