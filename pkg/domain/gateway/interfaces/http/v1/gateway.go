package v1

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xyedo/blindate/pkg/common/constant"
	"github.com/xyedo/blindate/pkg/common/event"
	"github.com/xyedo/blindate/pkg/domain/gateway"
	gatewayEntities "github.com/xyedo/blindate/pkg/domain/gateway/entities"
	"github.com/xyedo/blindate/pkg/infrastructure"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 20 * time.Second
	pingPeriod = 10 * time.Second
)

func New(config infrastructure.Config, gateway gateway.Session) *gatewayH {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			allowedOrigins := []string{"https://blindate.com"}
			if os.Getenv("env") == "development" {
				allowedOrigins = append(allowedOrigins, "http://localhost:3000")
			}
			reqOrigin := r.Header.Get("Origin")
			for _, allallowedOrigin := range allowedOrigins {
				if allallowedOrigin == reqOrigin {
					return true
				}
			}
			return false
		},
	}
	return &gatewayH{
		config:         config,
		upgrader:       &upgrader,
		session:        gateway,
		gatewayChannel: make(chan requestPayload, 10),
	}
}

type gatewayH struct {
	config         infrastructure.Config
	upgrader       *websocket.Upgrader
	session        gateway.Session
	gatewayChannel chan requestPayload
}

func (g *gatewayH) Listen() {
	for {
		event := <-g.gatewayChannel
		switch event.Action {
		case "onTypingStart":
			g.onSimpleAction(event, "onTypingStart")
		case "onTypingStop":
			g.onSimpleAction(event, "onTypingStop")
		case "onSendingVoiceStart":
			g.onSimpleAction(event, "onSendingVoiceStart")
		case "onSendingVoiceStop":
			g.onSimpleAction(event, "onSendingVoiceStop")
		case "onChoosingStickerStart":
			g.onSimpleAction(event, "onChoosingStickerStart")
		case "onChoosingStickerStop":
			g.onSimpleAction(event, "onChoosingStickerStop")
		case "onLeaving":
			g.cleanUp(event.UserId)

		}
	}
}

func (g *gatewayH) wsHandler(c *gin.Context) {
	userId := c.GetString(constant.KeyRequestUserId)
	wsConnection, err := g.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		log.Println(reflect.TypeOf(err))
		return
	}

	event.UserConnectionChange.Trigger(event.ConnectionPayload{
		UserId: userId,
		Online: true,
	})

	conn := gatewayEntities.Conn{Conn: wsConnection}
	g.session.SetUserSocket(userId, conn)

	go g.handleRequestPayload(userId, conn)
	go g.pingTicker(conn)

}

func (g *gatewayH) handleRequestPayload(id string, conn gatewayEntities.Conn) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error panic", err)
		}

		g.cleanUp(id)

	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))

	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	var payload requestPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}
		payload.UserId = id

		g.gatewayChannel <- payload
	}
}

func (ws *gatewayH) pingTicker(conn gatewayEntities.Conn) {
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

func (g *gatewayH) cleanUp(id string) {
	conn, ok := g.session.GetUserSocket(id)
	if !ok {
		return
	}
	conn.Close()

	g.session.DeleteUserSocket(id)

	event.UserConnectionChange.Trigger(event.ConnectionPayload{
		UserId: id,
		Online: false,
	})

}

func (g *gatewayH) onSimpleAction(event requestPayload, action string) {

	// convId := event.Payload
	// match, err := g.MatchSvc.GetMatchById(convId)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// sendToConversation(match.RequestFrom, convId)
	// sendToConversation(match.RequestTo, convId)
}
