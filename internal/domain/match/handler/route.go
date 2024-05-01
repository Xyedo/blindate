package handler

import "github.com/labstack/echo/v4"

func (h *Match) Route(e *echo.Group) {
	matchs := e.Group("/matchs")

	matchs.POST("", h.postCreateNewCandidateMatch)
	matchs.GET("", h.getIndexMatchs)

	matchs.GET("/:matchId", h.getMatchById)
	matchs.PUT("/:matchId/request-transition", h.putTransitionRequestStatus)
}
