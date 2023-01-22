package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

type locationSvc interface {
	CreateNewLocation(location *domain.Location) error
	UpdateLocation(userId string, changeLat, changeLng *string) error
	GetLocation(id string) (domain.Location, error)
}

func NewLocation(locationService locationSvc) *Location {
	return &Location{
		locationService: locationService,
	}
}

type Location struct {
	locationService locationSvc
}

func (l *Location) postLocationByUserIdHandler(c *gin.Context) {
	userId := c.GetString(keyUserId)
	var input struct {
		Lat string `json:"lat" binding:"required,latitude"`
		Lng string `json:"lng" binding:"required,longitude"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"lat": "required and must be valid lat geometry",
			"lng": "required and must be valid lng geometry",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = l.locationService.CreateNewLocation(&domain.Location{UserId: userId, Lat: input.Lat, Lng: input.Lng})
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "location created",
	})
}
func (l *Location) getLocationByUserIdHandler(c *gin.Context) {
	userId := c.GetString(keyUserId)
	location, err := l.locationService.GetLocation(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"location": location,
		},
	})
}

func (l *Location) patchLocationByUserIdHandler(c *gin.Context) {
	userId := c.GetString(keyUserId)
	var input struct {
		Lat *string `json:"lat" binding:"omitempty,latitude"`
		Lng *string `json:"lng" binding:"omitempty,longitude"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"lat": "must be valid lat geometry",
			"lng": "must be valid lng geometry",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = l.locationService.UpdateLocation(userId, input.Lat, input.Lng)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "location updated",
	})
}
