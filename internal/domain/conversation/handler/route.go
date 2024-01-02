package handler

import "github.com/labstack/echo/v4"

func Route(e *echo.Group) {
	conversations := e.Group("/conversations")

	conversations.GET("", getIndexConversations)

	conversations.GET("/:convId", getIndexChats)
}
