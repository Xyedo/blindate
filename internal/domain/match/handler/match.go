package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xyedo/blindate/internal/domain/match/dtos"
	"github.com/xyedo/blindate/internal/domain/match/usecase"
	"github.com/xyedo/blindate/internal/infrastructure/auth"
	"github.com/xyedo/blindate/pkg/pagination"
)

func postCreateNewCandidateMatch(c echo.Context) error {
	ctx := c.Request().Context()
	requestId := ctx.Value(auth.RequestId).(string)

	err := usecase.CreateCandidateMatch(ctx, requestId)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}

func getIndexMatchs(c echo.Context) error {
	ctx := c.Request().Context()

	var queryParams dtos.IndexMatchsQueryParams

	err := c.Bind(&queryParams)
	if err != nil {
		return err
	}

	err = queryParams.Mod().Validate()
	if err != nil {
		return err
	}

	requestId := ctx.Value(auth.RequestId).(string)
	matchedUsers, hasNext, err := usecase.IndexMatchCandidate(ctx, requestId, pagination.Pagination(queryParams))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK,
		dtos.NewIndexMatchResponse(
			hasNext,
			pagination.Pagination(queryParams),
			matchedUsers,
		),
	)
}

func putTransitionRequestStatus(c echo.Context) error {
	ctx := c.Request().Context()

	var request dtos.PutTransitionRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	requestId := ctx.Value(auth.RequestId).(string)

	return usecase.TransitionRequestStatus(ctx, requestId, c.Param("matchId"), request.Swipe)
}
