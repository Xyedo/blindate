package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

func NewLocation(locationService service.Location) *location {
	return &location{
		locationService: locationService,
	}
}

type location struct {
	locationService service.Location
}

func (l *location) postLocationByUserIdHandler(c *gin.Context) {
	userId := c.GetString("userId")
	var input struct {
		Lat string `json:"lat" binding:"required,latitude"`
		Lng string `json:"lng" binding:"required,longitude"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Lat": "required and must be valid lat geometry",
			"Lng": "required and must be valid lng geometry",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	err = l.locationService.CreateNewLocation(&domain.Location{UserId: userId, Lat: input.Lat, Lng: input.Lng})
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefUserIdField) {
			errorResourceNotFound(c, "user id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "location created",
	})
}
func (l *location) getLocationByUserIdHandler(c *gin.Context) {
	userId := c.GetString("userId")
	location, err := l.locationService.GetLocation(userId)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "user id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"location": location,
		},
	})
}

func (l *location) patchLocationByUserIdHandler(c *gin.Context) {
	userId := c.GetString("userId")
	var input struct {
		Lat *string `json:"lat" binding:"omitempty,latitude"`
		Lng *string `json:"lng" binding:"omitempty,longitude"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Lat": "must be valid lat geometry",
			"Lng": "must be valid lng geometry",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	location, err := l.locationService.GetLocation(userId)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "user id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	if input.Lat != nil {
		location.Lat = *input.Lat
	}
	if input.Lng != nil {
		location.Lng = *input.Lng
	}

	err = l.locationService.UpdateLocation(location)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "user id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "location updated",
	})
}
