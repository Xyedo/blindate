package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func (h *interestH) postTravelsHandler(c *gin.Context) {
	var request postTravelsRequest
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

	travelIds, err := h.interestUC.CreateTravelsByInterestId(
		c.GetString(constant.KeyInterestId),
		request.Travels,
	)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"travels": travelIds,
		},
	})
}

func (h *interestH) patchTravelsHandler(c *gin.Context) {
	var request patchTravelsRequest
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

	travelsDTO := make([]interestDTOs.Travel, 0, len(request.Travels))
	for _, travel := range request.Travels {
		travelsDTO = append(travelsDTO, interestDTOs.Travel(travel))
	}
	err = h.interestUC.UpdateTravels(travelsDTO)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "update travels success",
	})
}

func (h *interestH) deleteTravelsHandler(c *gin.Context) {
	var request deleteTravelsRequest
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

	err = h.interestUC.DeleteTravelsByIDs(request.IDs)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
