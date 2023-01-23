package gateway

import (
	"log"

	"github.com/xyedo/blindate/pkg/applications/service"
	websocketEntity "github.com/xyedo/blindate/pkg/domain/ws"
)

type Deps struct {
	Ws         *service.Ws
	ChatSvc    *service.Chat
	MatchSvc   *service.Match
	OnlinceSvc *service.Online
}

func (d *Deps) ListenToWsChan() {
	for {
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
		case "onChoosingStickerStart":
			d.OnSimpleAction(event, "onChoosingStickerStart")
		case "onChoosingStickerStop":
			d.OnSimpleAction(event, "onChoosingStickerStop")
		case "onLeaving":
			d.OnLeaving(event)
		}
	}

}
func (d *Deps) OnLeaving(event websocketEntity.Payload) {
	userId, ok := d.Ws.Clients.Get(event.Conn)
	if !ok {
		return
	}
	d.removingUser(event.Conn, userId)

}

func (d *Deps) OnSimpleAction(event websocketEntity.Payload, action string) {
	sendToConversation := func(toUserId, convId string) {
		socket, ok := d.Ws.ReverseClient.Get(toUserId)
		if !ok {
			return
		}
		err := socket.WriteJSON(websocketEntity.Response{
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
func (d *Deps) removingUser(socket websocketEntity.Conn, userId string) {
	_ = socket.Close()
	d.Ws.Clients.Delete(socket)
	d.Ws.ReverseClient.Delete(userId)
	d.OnlinceSvc.PutOnline(userId, false)
}
