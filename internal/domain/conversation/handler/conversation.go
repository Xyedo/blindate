package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xyedo/blindate/internal/domain/conversation/dtos"
	"github.com/xyedo/blindate/internal/domain/conversation/usecase"
	"github.com/xyedo/blindate/internal/infrastructure/auth"
	"github.com/xyedo/blindate/pkg/pagination"
)

func getIndexConversations(c echo.Context) error {
	ctx := c.Request().Context()

	var queryParams dtos.IndexConversationQueryParams

	err := c.Bind(&queryParams)
	if err != nil {
		return err
	}

	err = queryParams.Mod().Validate()
	if err != nil {
		return err
	}
	requestId := ctx.Value(auth.RequestId).(string)
	convs, hasNext, err := usecase.IndexConversation(ctx, requestId, queryParams.Page, queryParams.Limit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK,
		dtos.NewIndexConversationResponse(
			hasNext,
			pagination.Pagination(queryParams),
			convs,
		),
	)
}

func getIndexChats(c echo.Context) error {
	ctx := c.Request().Context()

	var queryParams dtos.IndexChatQueryParams
	err := c.Bind(&queryParams)
	if err != nil {
		return err
	}

	err = queryParams.Mod().Validate()
	if err != nil {
		return err
	}
	requestId := ctx.Value(auth.RequestId).(string)
	convs, hasNext, hasPrev, err := usecase.IndexChatByConversationId(ctx,
		queryParams.ToEntity(requestId, c.Param("convId")),
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK,
		dtos.NewIndexChatResponse(dtos.IndexChatPayload{
			HasNext: hasNext,
			HasPrev: hasPrev,
			Conv:    convs,
		}),
	)
}
