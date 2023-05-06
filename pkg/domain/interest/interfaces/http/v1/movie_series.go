package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func (h *interestH) postMovieSeriesHandler(c *gin.Context) {
	var request postMovieSeriesRequest

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

	movieSerieIds, err := h.interestUC.CreateMovieSeriesByInterestId(
		c.GetString(constant.KeyInterestId),
		request.MovieSeries,
	)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"movie_series": movieSerieIds,
		},
	})
}

func (h *interestH) patchMovieSeriesHandler(c *gin.Context) {
	var request patchMovieSeriesRequest
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

	movieSeriesDTO := make([]interestDTOs.MovieSerie, 0, len(request.MovieSeries))
	for _, movieSerie := range request.MovieSeries {
		movieSeriesDTO = append(movieSeriesDTO, interestDTOs.MovieSerie(movieSerie))
	}
	err = h.interestUC.UpdateMovieSeriesByInterestId(
		c.GetString(constant.KeyInterestId),
		movieSeriesDTO,
	)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "update movie_series success",
	})
}

func (h *interestH) deleteMovieSeriesHandler(c *gin.Context) {
	var request deleteMovieSeriesRequest
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

	err = h.interestUC.DeleteMovieSeriesByInterestId(
		c.GetString(constant.KeyInterestId),
		request.IDs,
	)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
