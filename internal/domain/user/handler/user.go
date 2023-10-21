package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/xyedo/blindate/internal/domain/user/usecase"
)

func handleClrekWebhook(c echo.Context) error {
	var request struct {
		Data struct {
			Id string `json:"id"`
		} `json:"data"`
		Type string `json:"type"`
	}

	err := c.Bind(&request)
	if err != nil {
		return err
	}
	switch request.Type {
	case "user.created":
		return usecase.RegisterUser(c.Request().Context(), request.Data.Id)
	}
	return nil
}
