package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	"github.com/xyedo/blindate/pkg/domain/interest"
	"github.com/xyedo/blindate/pkg/infrastructure"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func New(config infrastructure.Config, interestUsecase interest.Usecase) *interestH {
	return &interestH{
		config:     config,
		interestUC: interestUsecase,
	}
}

type interestH struct {
	config     infrastructure.Config
	interestUC interest.Usecase
}

func (h *interestH) getInterestDetailHandler(c *gin.Context) {
	interestDetail, err := h.interestUC.GetById(
		c.GetString(constant.KeyRequestUserId),
	)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}
	hobbies := make([]hobbie, 0, len(interestDetail.Hobbies))
	for _, hobbieDTO := range interestDetail.Hobbies {
		hobbies = append(hobbies, hobbie(hobbieDTO))
	}

	movieSeries := make([]movieSerie, 0, len(interestDetail.MovieSeries))
	for _, movieSerieDTO := range interestDetail.MovieSeries {
		movieSeries = append(movieSeries, movieSerie(movieSerieDTO))
	}

	travels := make([]travel, 0, len(interestDetail.Travels))
	for _, travelDTO := range interestDetail.Travels {
		travels = append(travels, travel(travelDTO))
	}

	sports := make([]sport, 0, len(interestDetail.Sports))
	for _, sportDTO := range interestDetail.Sports {
		sports = append(sports, sport(sportDTO))
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": getInterestDetailResponse{
			bio: bio{
				Bio:       interestDetail.Bio.Bio,
				CreatedAt: interestDetail.CreatedAt,
				UpdatedAt: interestDetail.UpdatedAt,
			},
			Hobbies:     hobbies,
			MovieSeries: movieSeries,
			Travels:     travels,
			Sports:      sports,
		},
	})
}
