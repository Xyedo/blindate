package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/internal/security"
	"github.com/xyedo/blindate/pkg/common/constant"
	httpMiddleware "github.com/xyedo/blindate/pkg/infrastructure/http/middleware"
)

func (h *interestH) Handler(globalRoute *gin.RouterGroup, jwt *security.Jwt) {
	interests := globalRoute.Group("/users/:userId/interests", httpMiddleware.AuthToken(jwt), httpMiddleware.ValidateUserId())

	interests.GET("/", h.getInterestDetailHandler)
	interests.POST("/bio", h.postBioHandler)
	interests.PATCH("/bio", h.patchBioHandler)

	interest := interests.Group("/:interestId", validateInterest())
	{
		interest.POST("/hobbies", h.postHobbiesHandler)
		interest.PATCH("/hobbies", h.patchHobbiesHandler)
		interest.DELETE("/hobbies", h.deleteHobbiesHandler)

		interest.POST("/movie-series", h.postMovieSeriesHandler)
		interest.PATCH("/movie-series", h.patchMovieSeriesHandler)
		interest.DELETE("/movie-series", h.deleteMovieSeriesHandler)
	}

}

func validateInterest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var url struct {
			InterestId string `uri:"interestId" binding:"required,uuid"`
		}
		err := c.ShouldBindUri(&url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  "fail",
				"message": "required,must have uuid in uri!",
			})
			return
		}
		c.Set(constant.KeyInterestId, url.InterestId)
		c.Next()
	}
}
