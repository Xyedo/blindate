package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

type interestSvc interface {
	GetInterest(userId string) (domain.Interest, error)

	CreateNewBio(intr *domain.Bio) error
	GetBio(userId string) (domain.Bio, error)
	PutBio(bio domain.Bio) error

	CreateNewHobbies(interestId string, hobbies []string) error
	PutHobbies(interestId string, hobbies []domain.Hobbie) error
	DeleteHobbies(interestId string, ids []string) error

	CreateNewMovieSeries(interestId string, movieSeries []string) error
	PutMovieSeries(interestId string, movieSeries []domain.MovieSerie) error
	DeleteMovieSeries(interestId string, ids []string) error

	CreateNewTraveling(interestId string, travels []string) error
	PutTraveling(interestId string, travels []domain.Travel) error
	DeleteTravels(interestId string, ids []string) error

	CreateNewSports(interestId string, sports []string) error
	PutSports(interestId string, sports []domain.Sport) error
	DeleteSports(interestId string, ids []string) error
}

func NewInterest(interestSvc interestSvc) *Interest {
	return &Interest{
		interestSvc: interestSvc,
	}
}

type Interest struct {
	interestSvc interestSvc
}

func (i *Interest) getInterestHandler(c *gin.Context) {
	userId := c.GetString("userId")
	intr, err := i.interestSvc.GetInterest(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"interest": intr,
		},
	})
}

func (i *Interest) postInterestBioHandler(c *gin.Context) {
	var input struct {
		Bio *string `json:"bio" binding:"required,max=300"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"bio": "at least an empty string and maximal character length is less than 300",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	userId := c.GetString("userId")
	bio := domain.Bio{
		UserId: userId,
		Bio:    *input.Bio,
	}
	err = i.interestSvc.CreateNewBio(&bio)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "interest bio created",
		"data": gin.H{
			"interestId": bio.Id,
		},
	})

}
func (i *Interest) putInterestBioHandler(c *gin.Context) {
	userId := c.GetString("userId")
	bio, err := i.interestSvc.GetBio(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	var input struct {
		Bio *string `json:"bio" binding:"required,max=300"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"bio": "required, maximal character length is less than 300",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	changedBio := strings.TrimSpace(*input.Bio)
	if bio.Bio == changedBio {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "nothing changed",
		})
		return
	}
	bio.Bio = changedBio
	err = i.interestSvc.PutBio(bio)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "bio successfully changed",
	})

}

func (i *Interest) postInterestHobbiesHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Hobbies []string `json:"hobbies" binding:"required,max=10,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"hobbies": "each hobbies must be unique, less than 10 and has more than 2 and less than 50 character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.CreateNewHobbies(interestId, input.Hobbies)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"hobbies": input.Hobbies,
		},
	})

}

func (i *Interest) putInterestHobbiesHandler(c *gin.Context) {
	var input struct {
		Hobbies []domain.Hobbie `json:"hobbies" binding:"required,max=10,unique=Hobbie"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"hobbies": "each hobbies must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new hobbies",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.PutHobbies(interestId, input.Hobbies)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"hobbies": input.Hobbies,
		},
	})

}
func (i *Interest) deleteInterestHobbiesHandler(c *gin.Context) {
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"ids": "each ids must be unique and uuid character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	interestId := c.GetString("interestId")
	err = i.interestSvc.DeleteHobbies(interestId, input.Ids)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "ids is successfully deleted",
	})

}

func (i *Interest) postInterestMovieSeriesHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		MovieSeries []string `json:"movieSeries" binding:"required,max=10,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"movieSeries": "each movieSeries must be unique, less than 10 and has more than 2 and less than 50 character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = i.interestSvc.CreateNewMovieSeries(interestId, input.MovieSeries)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"movieSeries": input.MovieSeries,
		},
	})
}
func (i *Interest) putInterestMovieSeriesHandler(c *gin.Context) {
	var input struct {
		MovieSeries []domain.MovieSerie `json:"movieSeries" binding:"required,max=10,unique=MovieSerie"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"movieSeries": "each movieSeries must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new movieSeries",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.PutMovieSeries(interestId, input.MovieSeries)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"movieSeries": input.MovieSeries,
		},
	})

}
func (i *Interest) deleteInterestMovieSeriesHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"ids": "each ids must be unique and uuid character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = i.interestSvc.DeleteMovieSeries(interestId, input.Ids)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "movieSeries is successfully deleted",
	})

}

func (i *Interest) postInterestTravelingHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Travels []string `json:"travels" binding:"required,max=10,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"travels": "each travels must be unique, less than 10 and has more than 2 and less than 50 character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.CreateNewTraveling(interestId, input.Travels)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"travels": input.Travels,
		},
	})
}

func (i *Interest) putInterestTravelingHandler(c *gin.Context) {
	var input struct {
		Travels []domain.Travel `json:"travels" binding:"required,max=10,unique=Travel"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"travels": "each travels must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new travel.",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.PutTraveling(interestId, input.Travels)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"travels": input.Travels,
		},
	})
}
func (i *Interest) deleteInterestTravelingHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"ids": "each ids must be unique and uuid character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = i.interestSvc.DeleteTravels(interestId, input.Ids)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "travels is successfully deleted",
	})

}

func (i *Interest) postInterestSportHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Sports []string `json:"sports" binding:"required,max=10,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"sports": "each sports must be unique, less than 10 and has more than 2 and less than 50 character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.CreateNewSports(interestId, input.Sports)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"sports": input.Sports,
		},
	})
}
func (i *Interest) putInterestSportHandler(c *gin.Context) {
	var input struct {
		Sports []domain.Sport `json:"sports" binding:"required,max=10,unique=Sport"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"sports": "each sports must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new sports.",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.PutSports(interestId, input.Sports)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"sports": input.Sports,
		},
	})
}
func (i *Interest) deleteInterestSportHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"ids": "each ids must be unique and uuid character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = i.interestSvc.DeleteSports(interestId, input.Ids)
	if err != nil {
		jsonHandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "sports is successfully deleted",
	})
}
