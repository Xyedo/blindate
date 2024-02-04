package handler

import "github.com/labstack/echo/v4"

func Route(e *echo.Group) {
	matchs := e.Group("/matchs")

	matchs.POST("", postCreateNewCandidateMatch)
	matchs.GET("", getIndexMatchs)

	matchs.GET("/:matchId", getMatchById)
	matchs.PUT("/:matchId/request-transition", putTransitionRequestStatus)
}
