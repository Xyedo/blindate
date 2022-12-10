package gateway

import (
	"log"

	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

type WsDeps struct {
	Ws         *service.Ws
	ChatSvc    *service.Chat
	MatchSvc   *service.Match
	OnlinceSvc *service.Online
}

func (d *WsDeps) ListenToWsChan() {
	event := <-d.Ws.WsChan
	switch event.Action {
	case "onTypingStart":
		d.OnSimpleAction(event, "onTypingStart")
	case "onTypingStop":
		d.OnSimpleAction(event, "onTypingStop")
	case "onSendingVoiceStart":
		d.OnSimpleAction(event, "onSendingVoiceStart")
	case "onSendingVoiceStop":
		d.OnSimpleAction(event, "onSendingVoiceStop")
	case "onLeaving":

	}
}
func (d *WsDeps) OnLeaving(event domain.WsPayload) {
	userId, ok := d.Ws.Clients[event.Conn]
	if !ok {
		return
	}
	d.removingUser(event.Conn, userId)

}
func (d *WsDeps) removingUser(socket domain.WsConn, userId string) {
	_ = socket.Close()
	delete(d.Ws.Clients, socket)
	delete(d.Ws.ReverseClient, userId)
	d.OnlinceSvc.PutOnline(userId, false)
}
func (d *WsDeps) OnSimpleAction(event domain.WsPayload, action string) {
	sendToConversation := func(toUserId, convId string) {
		socket, ok := d.Ws.ReverseClient[toUserId]
		if !ok {
			return
		}
		err := socket.WriteJSON(domain.WsResponse{
			Action: action,
			Data: map[string]any{
				"convId": convId,
			},
		})
		if err != nil {
			log.Println("websocket Err", err)
			d.removingUser(socket, toUserId)
		}
	}
	convId := event.Payload
	match, err := d.MatchSvc.GetMatchById(convId)
	if err != nil {
		log.Println(err)
		return
	}
	sendToConversation(match.RequestFrom, convId)
	sendToConversation(match.RequestTo, convId)
}
