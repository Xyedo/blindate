package external

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xyedo/blindate/internal/domain/user/dtos"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/domain/user/usecase"
	"github.com/xyedo/blindate/internal/infrastructure/auth"
)

func postBasicInfo(c echo.Context) error {
	var request dtos.PostBasicInfoRequest

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	requestId := ctx.Value(auth.RequestId).(string)

	returnedId, err := usecase.CreateBasicInfo(ctx, requestId, entities.CreateBasicInfo{
		Gender:           request.Gender,
		FromLoc:          request.FromLoc,
		Height:           request.Height,
		EducationLevel:   request.EducationLevel,
		Drinking:         request.Drinking,
		Smoking:          request.Smoking,
		RelationshipPref: request.RelationshipPref,
		LookingFor:       request.LookingFor,
		Zodiac:           request.Zodiac,
		Kids:             request.Kids,
		Work:             request.Work,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data": map[string]any{
			"id": returnedId,
		},
	})

}

func getBasicInfoById(c echo.Context) error {
	ctx := c.Request().Context()
	requestId := ctx.Value(auth.RequestId).(string)

	returnedBasicInfo, err := usecase.GetBasicInfoById(ctx, requestId, c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data": returnedBasicInfo,
	})
}

func patchBasicInfoById(c echo.Context) error {
	var request dtos.PatchBasicInfoRequest

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = request.Validate()
	if err != nil {
		return err
	}
	ctx := c.Request().Context()
	requestId := ctx.Value(auth.RequestId).(string)

	return usecase.UpdateBasicInfoById(ctx, requestId, entities.UpdateBasicInfo{})

}
