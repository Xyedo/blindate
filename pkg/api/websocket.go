package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

func NewWs(wsSvc *service.Ws, onlineSvc onlineSvc) *Ws {
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
	return &Ws{
		wsSvc:     wsSvc,
		upgrader:  &upgrader,
		onlineSvc: onlineSvc,
	}
}

type Ws struct {
	wsSvc     *service.Ws
	upgrader  *websocket.Upgrader
	onlineSvc onlineSvc
}

func (ws *Ws) wsEndPoint(c *gin.Context) {
	userId := c.GetString("userId")
	wsConn, err := ws.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	ws.onlineSvc.PutOnline(userId, true)
	conn := domain.WsConn{Conn: wsConn}
	ws.wsSvc.Clients[conn] = userId
	ws.wsSvc.ReverseClient[userId] = conn

	go ws.wsSvc.ListenForWsPayload(&conn)
}
