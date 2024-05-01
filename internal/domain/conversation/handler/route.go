package handler

import "github.com/labstack/echo/v4"

func (h *Conversation) Route(e *echo.Group) {
	conversations := e.Group("/conversations")

	conversations.GET("", h.getIndexConversations)

	conversations.GET("/:convId", h.getIndexChats)
}
