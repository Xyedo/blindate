package service

import (
	"fmt"
	"log"
	"time"

	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/domain/entity"
	"github.com/xyedo/blindate/pkg/event"
)

type EventDeps struct {
	UserSvc  *User
	ConvSvc  *Conversation
	MatchSvc *Match
	Online   *Online
	Ws       *Ws
}

func (d *EventDeps) HandleSeenAtevent(payload event.ChatSeenPayload) {
	sendToChat := func(reqUserId string, updatedChatIds []string) {
		socket, ok := d.Ws.ReverseClient.Get(reqUserId)
		if !ok {
			return
		}
		var response domain.WsResponse
		response.Action = "update.chat.seenAt"
		response.Data = map[string]any{
			"seenChatIds": updatedChatIds,
		}
		socket.SetWriteDeadline(time.Now().Add(writeWait))
		err := socket.WriteJSON(response)
		if err != nil {
			log.Println("webscoket err", err)
			d.removingUser(socket, reqUserId)
		}
	}
	sendToChat(payload.RequestFrom, payload.SeenChatIds)
	sendToChat(payload.RequestFrom, payload.SeenChatIds)
}

func (d *EventDeps) HandleProfileUpdateEvent(payload event.ProfileUpdatedPayload) {
	sendToConversation := func(updatedUser domain.User, userId, convUserId string) {
		socket, ok := d.Ws.ReverseClient.Get(convUserId)
		if !ok {
			return
		}
		var response domain.WsResponse
		response.Action = "update.conversation.profile"
		response.Data = map[string]any{"updatedUser": updatedUser}
		socket.SetWriteDeadline(time.Now().Add(writeWait))
		err := socket.WriteJSON(response)
		if err != nil {
			log.Println("webscoket err", err)
			d.removingUser(socket, convUserId)
		}
	}

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
		sendToConversation(updatedUser, payload.UserId, conv.FromUser.ID)
		sendToConversation(updatedUser, payload.UserId, conv.ToUser.ID)
	}
}

func (d *EventDeps) HandleRevealUpdateEvent(payload event.MatchRevealedPayload) {
	sendToMatch := func(reqUserId string, match entity.Match) {
		conn, ok := d.Ws.ReverseClient.Get(reqUserId)
		if !ok {
			return
		}
		response := domain.WsResponse{
			Action: fmt.Sprintf("reveal.%s", payload.MatchStatus),
			Data: map[string]any{
				"match": match,
			},
		}
		conn.SetWriteDeadline(time.Now().Add(writeWait))
		err := conn.WriteJSON(response)
		if err != nil {
			log.Println("webscoket err", err)
			d.removingUser(conn, reqUserId)
		}
	}
	matchEntity, err := d.MatchSvc.GetMatchById(payload.MatchId)
	if err != nil {
		log.Println(err)
		return
	}
	sendToMatch(matchEntity.RequestFrom, matchEntity)
	sendToMatch(matchEntity.RequestTo, matchEntity)

}
func (d *EventDeps) HandleCreateChatEvent(payload event.ChatCreatedPayload) {
	sendResponse := func(socketConn domain.WsConn, userId string, chat []domain.Chat, conv domain.Conversation) {
		socketConn.SetWriteDeadline(time.Now().Add(writeWait))
		err := socketConn.WriteJSON(domain.WsResponse{
			Action: "OnMessage",
			Data: map[string]any{
				"chats": chat,
				"conv":  conv,
			},
		})
		if err != nil {
			log.Println("webscoket err", err)
			d.removingUser(socketConn, userId)
		}
	}
	conv, err := d.ConvSvc.FindConversationById(payload.ConvId)
	if err != nil {
		log.Println(err)
		return
	}
	if authorSoc, ok := d.Ws.ReverseClient.Get(conv.FromUser.ID); ok {
		sendResponse(authorSoc, conv.FromUser.ID, payload.Chat, conv)
	}
	if recipientSoc, ok := d.Ws.ReverseClient.Get(conv.ToUser.ID); ok {
		sendResponse(recipientSoc, conv.ToUser.ID, payload.Chat, conv)
	}
}
func (d *EventDeps) removingUser(socket domain.WsConn, reqUserId string) {
	_ = socket.Close()
	d.Ws.Clients.Delete(socket)
	d.Ws.ReverseClient.Delete(reqUserId)
	_ = d.Online.PutOnline(reqUserId, false)
}
