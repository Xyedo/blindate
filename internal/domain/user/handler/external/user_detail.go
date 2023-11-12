package external

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xyedo/blindate/internal/domain/user/dtos"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/domain/user/usecase"
	"github.com/xyedo/blindate/internal/infrastructure/auth"
)

func postUserDetailHandler(c echo.Context) error {
	var request dtos.PostUserDetailRequest

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	returnedId, err := usecase.CreateUserDetail(c.Request().Context(), c.Param("id"), entities.CreateUserDetail{
		Gender:           request.Gender,
		Geog:             entities.Geography(request.Location),
		Bio:              request.Bio,
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

func getUserDetailByIdHandler(c echo.Context) error {
	ctx := c.Request().Context()
	requestId := ctx.Value(auth.RequestId).(string)

	userDetail, err := usecase.GetUserDetail(ctx, requestId, c.Param("id"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data": dtos.GetUserDetailResponse(userDetail),
	})
}

func patchUserDetailByIdHandler(c echo.Context) error {
	var request dtos.PatchUserDetailRequest

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

	return usecase.UpdateUserDetailById(ctx, requestId, request.ToEntity())

}
