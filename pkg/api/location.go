package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

type locationSvc interface {
	CreateNewLocation(location *domain.Location) error
	UpdateLocation(location *domain.Location) error
	GetLocation(id string) (*domain.Location, error)
}

func NewLocation(locationService locationSvc) *location {
	return &location{
		locationService: locationService,
	}
}

type location struct {
	locationService locationSvc
}

func (l *location) postLocationByUserIdHandler(c *gin.Context) {
	userId := c.GetString("userId")
	var input struct {
		Lat string `json:"lat" binding:"required,latitude"`
		Lng string `json:"lng" binding:"required,longitude"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Lat": "required and must be valid lat geometry",
			"Lng": "required and must be valid lng geometry",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = l.locationService.CreateNewLocation(&domain.Location{UserId: userId, Lat: input.Lat, Lng: input.Lng})
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefUserIdField) {
			errNotFoundResp(c, "user id is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainUserId) {
			errUnprocessableEntityResp(c, "location with this user id is already created")
			return
		}
		errServerResp(c, err)
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
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "user id is not found")
			return
		}
		errServerResp(c, err)
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
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Lat": "must be valid lat geometry",
			"Lng": "must be valid lng geometry",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	location, err := l.locationService.GetLocation(userId)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "user id is not found")
			return
		}
		errServerResp(c, err)
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
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "user id is not found")
			return
		}
		errServerResp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "location updated",
	})
}
