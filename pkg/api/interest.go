package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

func NewInterest(interestSvc service.Interest) *interest {
	return &interest{
		interestSvc: interestSvc,
	}
}

type interest struct {
	interestSvc service.Interest
}

func (i *interest) getInterestHandler(c *gin.Context) {
	userId := c.GetString("userId")
	intr, err := i.interestSvc.GetInterest(userId)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {

			errorResourceNotFound(c, "userId is not match with our resource")
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"interest": intr,
		},
	})
}

func (i *interest) postInterestBioHandler(c *gin.Context) {
	var input struct {
		Bio string `json:"bio" binding:"omitempty,max=300"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}

		errMap := util.ReadValidationErr(err, map[string]string{
			"Bio": "maximal character length is less than 300",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	userId := c.GetString("userId")
	bio := domain.Bio{
		UserId: userId,
		Bio:    input.Bio,
	}

	err = i.interestSvc.CreateNewBio(&bio)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefUserIdField) {
			errorResourceNotFound(c, "userId is not match with our resource")
			return
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "interest created",
	})

}
func (i *interest) putInterestBioHandler(c *gin.Context) {
	userId := c.GetString("userId")
	bio, err := i.interestSvc.GetBio(userId)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "userId is not match with our resource")
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		errorServerResponse(c, err)
		return
	}
	var input struct {
		Bio *string `json:"bio" binding:"omitempty,max=300"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Bio": "maximal character length is less than 300",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	if input.Bio != nil {
		bio.Bio = *input.Bio
	} else {
		c.JSON(http.StatusNotModified, gin.H{
			"status":  "success",
			"message": "nothing changed",
		})
	}
	err = i.interestSvc.PutBio(bio)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "userId is not found in our db!")
			return
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "bio successfully changed",
	})

}

func (i *interest) postInterestHobbiesHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Hobbies []string `json:"hobbies" binding:"required,max=25,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Hobbies": "each hobbies must be unique and has more than 2 and less than 50 character",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	hobbies := make([]domain.Hobbie, 0, len(input.Hobbies))
	for _, hobbie := range input.Hobbies {
		hobbies = append(hobbies, domain.Hobbie{
			Hobbie: hobbie,
		})
	}
	err = i.interestSvc.CreateNewHobbies(interestId, hobbies)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errorResourceNotFound(c, "interestId is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"hobbies": hobbies,
		},
	})

}

func (i *interest) putInterestHobbiesHandler(c *gin.Context) {
	var input struct {
		Hobbies []domain.Hobbie `json:"hobbies" binding:"required,max=25,unique=Hobbie"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Hobbies": "each hobbies must be unique and has more than 2 and less than 50 character. Id must match or empty when its new hobbies.",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}

	err = i.interestSvc.PutHobbies(interestId, input.Hobbies)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errorResourceNotFound(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "all of the hobbies must be unique",
			})
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"hobbies": input.Hobbies,
		},
	})

}
func (i *interest) deleteInterestHobbiesHandler(c *gin.Context) {
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Ids": "each ids must be unique and uuid character",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	err = i.interestSvc.DeleteHobbies(input.Ids)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "every id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "ids is successfully deleted",
	})

}

func (i *interest) postInterestMovieSeriesHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		MovieSeries []string `json:"movieSeries" binding:"required,max=25,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"MovieSeries": "each movieSeries must be unique and has more than 2 and less than 50 character",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	movieSeries := make([]domain.MovieSerie, 0, len(input.MovieSeries))
	for _, movieSerie := range input.MovieSeries {
		movieSeries = append(movieSeries, domain.MovieSerie{
			MovieSerie: movieSerie,
		})
	}
	err = i.interestSvc.CreateNewMovieSeries(interestId, movieSeries)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errorResourceNotFound(c, "interestId is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"movieSeries": movieSeries,
		},
	})
}
func (i *interest) putInterestMovieSeriesHandler(c *gin.Context) {
	var input struct {
		MovieSeries []domain.MovieSerie `json:"movieSeries" binding:"required,max=25,unique=MovieSerie"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"MovieSeries": "each movieSeries must be unique and has more than 2 and less than 50 character. Id must match or empty when its new movieSeries.",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}

	err = i.interestSvc.PutMovieSeries(interestId, input.MovieSeries)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errorResourceNotFound(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "all of the movieSeries must be unique",
			})
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"movieSeries": input.MovieSeries,
		},
	})

}
func (i *interest) deleteInterestMovieSeriesHandler(c *gin.Context) {
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Ids": "each ids must be unique and uuid character",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	err = i.interestSvc.DeleteMovieSeries(input.Ids)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "every id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "movieSeries is successfully deleted",
	})

}

func (i *interest) postInterestTravelingHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Travels []string `json:"travels" binding:"required,max=25,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"MovieSeries": "each movieSeries must be unique and has more than 2 and less than 50 character",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	travels := make([]domain.Travel, 0, len(input.Travels))
	for _, travel := range input.Travels {
		travels = append(travels, domain.Travel{
			Travel: travel,
		})
	}
	err = i.interestSvc.CreateNewTraveling(interestId, travels)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errorResourceNotFound(c, "interestId is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"travels": travels,
		},
	})
}

func (i *interest) putInterestTravelingHandler(c *gin.Context) {
	var input struct {
		Travels []domain.Travel `json:"travels" binding:"required,max=25,unique=Travel"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Travels": "each travels must be unique and has more than 2 and less than 50 character. Id must match or empty when its new travel.",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}

	err = i.interestSvc.PutTraveling(interestId, input.Travels)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errorResourceNotFound(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "all of the travels must be unique",
			})
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"travels": input.Travels,
		},
	})
}
func (i *interest) deleteInterestTravelingHandler(c *gin.Context) {
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Ids": "each ids must be unique and uuid character",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	err = i.interestSvc.DeleteTravels(input.Ids)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "every id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "travels is successfully deleted",
	})

}

func (i *interest) postInterestSportHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Sports []string `json:"sports" binding:"required,max=25,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"MovieSeries": "each movieSeries must be unique and has more than 2 and less than 50 character",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	sports := make([]domain.Sport, 0, len(input.Sports))
	for _, sport := range input.Sports {
		sports = append(sports, domain.Sport{
			Sport: sport,
		})
	}
	err = i.interestSvc.CreateNewSports(interestId, sports)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errorResourceNotFound(c, "interestId is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"sports": sports,
		},
	})
}
func (i *interest) putInterestSportHandler(c *gin.Context) {
	var input struct {
		Sports []domain.Sport `json:"sports" binding:"required,max=25,unique=Sport"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Sports": "each sports must be unique and has more than 2 and less than 50 character. Id must match or empty when its new sport.",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}

	err = i.interestSvc.PutSports(interestId, input.Sports)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errorResourceNotFound(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "fail",
				"message": "all of the sports must be unique",
			})
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"travels": input.Sports,
		},
	})
}
func (i *interest) deleteInterestSportHandler(c *gin.Context) {
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Ids": "each ids must be unique and uuid character",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	err = i.interestSvc.DeleteSports(input.Ids)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "every id is not found")
			return
		}
		errorServerResponse(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "sports is successfully deleted",
	})
}
