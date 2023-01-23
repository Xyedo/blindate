package event

import (
	matchEntity "github.com/xyedo/blindate/pkg/domain/match/entities"
)

var MatchRevealed matchRevealed

type MatchRevealedPayload struct {
	MatchId     string
	MatchStatus matchEntity.Status
}

type matchRevealed struct {
	handlers []interface{ HandleRevealUpdateEvent(MatchRevealedPayload) }
}

func (m *matchRevealed) Register(handler interface{ HandleRevealUpdateEvent(MatchRevealedPayload) }) {
	m.handlers = append(m.handlers, handler)
}

func (m matchRevealed) Trigger(payload MatchRevealedPayload) {
	for _, handler := range m.handlers {
		go handler.HandleRevealUpdateEvent(payload)
	}
}
