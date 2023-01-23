package service

import (
	"fmt"
	"log"
	"time"

	"github.com/xyedo/blindate/pkg/domain/event"
	websocketEntity "github.com/xyedo/blindate/pkg/domain/ws"
)

type EventDeps struct {
	UserSvc  *User
	ConvSvc  *Conversation
	MatchSvc *Match
	Online   *Online
	Ws       *Ws
}

func (d *EventDeps) HandleSeenAtevent(payload event.ChatSeenPayload) {
	var response websocketEntity.Response
	response.Action = "update.chat.seenAt"
	response.Data = map[string]any{
		"seenChatIds": payload.SeenChatIds,
	}
	d.eventWriteJSON(payload.RequestFrom, response)
	d.eventWriteJSON(payload.RequestTo, response)
}

func (d *EventDeps) HandleProfileUpdateEvent(payload event.ProfileUpdatedPayload) {
	convs, err := d.ConvSvc.GetConversationByUserId(payload.UserId)
	if err != nil {
		log.Println(err)
		return
	}
	updatedUser, err := d.UserSvc.GetUserByIdWithSelectedProfPic(payload.UserId)
	if err != nil {
		log.Println(err)
		return
	}
	for _, conv := range convs {
		if conv.FromUser.ID == payload.UserId || conv.ToUser.ID == payload.UserId {
			continue
		}

		var response websocketEntity.Response
		response.Action = "update.conversation.profile"
		response.Data = map[string]any{"updatedUser": updatedUser}

		d.eventWriteJSON(conv.FromUser.ID, response)
		d.eventWriteJSON(conv.ToUser.ID, response)
	}
}

func (d *EventDeps) HandleRevealUpdateEvent(payload event.MatchRevealedPayload) {
	matchEntity, err := d.MatchSvc.GetMatchById(payload.MatchId)
	if err != nil {
		log.Println(err)
		return
	}
	response := websocketEntity.Response{
		Action: fmt.Sprintf("reveal.%s", payload.MatchStatus),
		Data: map[string]any{
			"match": matchEntity,
		},
	}
	d.eventWriteJSON(matchEntity.RequestFrom, response)
	d.eventWriteJSON(matchEntity.RequestTo, response)

}
func (d *EventDeps) HandleCreateChatEvent(payload event.ChatCreatedPayload) {
	conv, err := d.ConvSvc.FindConversationById(payload.ConvId)
	if err != nil {
		log.Println(err)
		return
	}
	resp := websocketEntity.Response{
		Action: "OnMessage",
		Data: map[string]any{
			"chats": payload.Chat,
			"conv":  conv,
		},
	}
	d.eventWriteJSON(conv.FromUser.ID, resp)
	d.eventWriteJSON(conv.ToUser.ID, resp)
}

func (d *EventDeps) eventWriteJSON(userId string, resp websocketEntity.Response) {
	socket, ok := d.Ws.ReverseClient.Get(userId)
	if !ok {
		return
	}
	socket.SetWriteDeadline(time.Now().Add(writeWait))
	err := socket.WriteJSON(resp)
	if err != nil {
		log.Println("webscoket err", err)
		socket.Close()
		d.Ws.Clients.Delete(socket)
		d.Ws.ReverseClient.Delete(userId)
		d.Online.PutOnline(userId, false)
	}
}
