package service

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xyedo/blindate/internal/rwmap"
	websocketEntity "github.com/xyedo/blindate/pkg/domain/ws"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 20 * time.Second
	pingPeriod = 10 * time.Second
)

func NewWs() *Ws {
	return &Ws{
		Clients:       rwmap.New[websocketEntity.Conn, string](),
		ReverseClient: rwmap.New[string, websocketEntity.Conn](),
		WsChan:        make(chan websocketEntity.Payload),
	}
}

type Ws struct {
	Clients       *rwmap.RwMap[websocketEntity.Conn, string]
	ReverseClient *rwmap.RwMap[string, websocketEntity.Conn]
	WsChan        chan websocketEntity.Payload
}

func (ws *Ws) ListenForWsPayload(conn *websocketEntity.Conn) {
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

	var payload websocketEntity.Payload
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

func (ws *Ws) PingTicker(conn *websocketEntity.Conn) {
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

func (ws *Ws) cleanUp(conn *websocketEntity.Conn) {
	conn.Close()

	userId, ok := ws.Clients.Get(*conn)
	if !ok {
		return
	}
	ws.Clients.Delete(*conn)
	ws.ReverseClient.Delete(userId)

}
