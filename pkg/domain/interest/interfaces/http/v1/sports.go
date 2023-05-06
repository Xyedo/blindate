package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func (h *interestH) postSportsHandler(c *gin.Context) {
	var request postSportsRequest
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

	sportIds, err := h.interestUC.CreateSportsByInterestId(
		c.GetString(constant.KeyInterestId),
		request.Sports,
	)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"sports": sportIds,
		},
	})
}

func (h *interestH) patchSportsHandler(c *gin.Context) {
	var request patchSportsRequest
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

	sportsDTO := make([]interestDTOs.Sport, 0, len(request.Sports))
	for _, sport := range request.Sports {
		sportsDTO = append(sportsDTO, interestDTOs.Sport(sport))
	}
	err = h.interestUC.UpdateSports(sportsDTO)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "update sports success",
	})
}
func (h *interestH) deleteSportsHandler(c *gin.Context) {
	var request deleteSportsRequest
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

	err = h.interestUC.DeleteSportsByIDs(request.IDs)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
