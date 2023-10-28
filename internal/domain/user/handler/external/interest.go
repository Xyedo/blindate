package external

import (
	"github.com/labstack/echo/v4"
	"github.com/xyedo/blindate/internal/domain/user/dtos"
	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/domain/user/usecase"
)

func postInterestHandler(c echo.Context) error {
	var request dtos.PostInterestRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = request.Mod().Validate()
	if err != nil {
		return err
	}

	return usecase.CreateInterest(
		c.Request().Context(),
		c.Param("id"),
		entities.CreateInterest(request),
	)

}

func patchInterestHandler(c echo.Context) error {
	var request dtos.PatchInterestRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = request.Mod().Validate()
	if err != nil {
		return err
	}

	return usecase.UpdateInterest(
		c.Request().Context(),
		c.Param("id"),
		request.ToEntity(),
	)

}

func postDeleteInterestHandler(c echo.Context) error {
	var request dtos.PostDeleteInterestRequest

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	return usecase.DeleteInterest(
		c.Request().Context(),
		c.Param("id"),
		entities.DeleteInterest(request),
	)
}
