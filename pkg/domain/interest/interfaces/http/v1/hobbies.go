package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func (h *interestH) postHobbiesHandler(c *gin.Context) {
	var request postHobbiesRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	err = request.mod().validate()
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	hobbieIds, err := h.interestUC.CreateHobbiesByInterestId(
		c.GetString(constant.KeyInterestId),
		request.Hobies,
	)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"hobbies": hobbieIds,
		},
	})
}

func (h *interestH) patchHobbiesHandler(c *gin.Context) {
	var request patchHobbiesRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	err = request.mod().validate()
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	hobbiesDTO := make([]interestDTOs.Hobbie, 0, len(request.Hobies))
	for _, hobbie := range request.Hobies {
		hobbiesDTO = append(hobbiesDTO, interestDTOs.Hobbie(hobbie))
	}
	err = h.interestUC.UpdateHobbiesByInterestId(
		c.GetString(constant.KeyInterestId),
		hobbiesDTO,
	)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "update hobbies success",
	})
}

func (h *interestH) deleteHobbiesHandler(c *gin.Context) {
	var request deleteHobbiesRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}
	err = request.mod().validate()
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	err = h.interestUC.DeleteHobbiesByInterestId(
		c.GetString(constant.KeyInterestId),
		request.IDs,
	)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
