package external

import (
	"github.com/labstack/echo/v4"
	"github.com/xyedo/blindate/internal/domain/user/dtos"
	"github.com/xyedo/blindate/internal/domain/user/entities"
)

func (u *User) postInterestHandler(c echo.Context) error {
	var request dtos.PostInterestRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = request.Mod().Validate()
	if err != nil {
		return err
	}

	return u.usecase.CreateInterest(
		c.Request().Context(),
		c.Param("id"),
		entities.CreateInterest(request),
	)

}

func (u *User) patchInterestHandler(c echo.Context) error {
	var request dtos.PatchInterestRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = request.Mod().Validate()
	if err != nil {
		return err
	}

	return u.usecase.UpdateInterest(
		c.Request().Context(),
		c.Param("id"),
		request.ToEntity(),
	)

}

func (u *User) postDeleteInterestHandler(c echo.Context) error {
	var request dtos.PostDeleteInterestRequest

	err := c.Bind(&request)
	if err != nil {
		return err
	}

	err = request.Validate()
	if err != nil {
		return err
	}

	return u.usecase.DeleteInterest(
		c.Request().Context(),
		c.Param("id"),
		entities.DeleteInterest(request),
	)
}
