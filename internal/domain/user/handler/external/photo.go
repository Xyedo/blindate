package external

import (
	"net/http"

	"github.com/labstack/echo/v4"
	apperror "github.com/xyedo/blindate/internal/common/app-error"
	"github.com/xyedo/blindate/internal/domain/user/usecase"
)

func putUserDetailPhotoHandler(c echo.Context) error {
	header, err := c.FormFile("photo")
	if err != nil {
		return err
	}
	if header == nil {
		return apperror.BadPayload(apperror.Payload{
			Message: "empty file",
		})
	}

	photoId, err := usecase.AddPhoto(c.Request().Context(), c.Param("id"), header)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"data": map[string]any{
			"photo_id": photoId,
		},
	})
}
