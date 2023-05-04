package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func (h *interestH) postBioHandler(c *gin.Context) {
	var request postBioRequest

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

	bioId, err := h.interestUC.CreateBio(interestDTOs.Bio{
		UserId: c.GetString(constant.KeyRequestUserId),
		Bio:    request.Bio,
	})
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"id": bioId,
		},
	})
}

func (h *interestH) patchBioHandler(c *gin.Context) {
	var request patchBioRequest

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

	err = h.interestUC.UpdateBio(interestDTOs.UpdateBio{
		UserId: c.GetString(constant.KeyRequestUserId),
		Bio:    request.Bio,
	})
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "patch bio success",
	})
}

