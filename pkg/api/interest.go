package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

type interestSvc interface {
	GetInterest(userId string) (*domain.Interest, error)

	CreateNewBio(intr *domain.Bio) error
	GetBio(userId string) (*domain.Bio, error)
	PutBio(bio *domain.Bio) error

	CreateNewHobbies(interestId string, hobbies []domain.Hobbie) error
	PutHobbies(interestId string, hobbies []domain.Hobbie) error
	DeleteHobbies(interestId string, ids []string) error

	CreateNewMovieSeries(interestId string, movieSeries []domain.MovieSerie) error
	PutMovieSeries(interestId string, movieSeries []domain.MovieSerie) error
	DeleteMovieSeries(interestId string, ids []string) error

	CreateNewTraveling(interestId string, travels []domain.Travel) error
	PutTraveling(interestId string, travels []domain.Travel) error
	DeleteTravels(interestId string, ids []string) error

	CreateNewSports(interestId string, sports []domain.Sport) error
	PutSports(interestId string, sports []domain.Sport) error
	DeleteSports(interestId string, ids []string) error
}

func NewInterest(interestSvc interestSvc) *interest {
	return &interest{
		interestSvc: interestSvc,
	}
}

type interest struct {
	interestSvc interestSvc
}

func (i *interest) getInterestHandler(c *gin.Context) {
	userId := c.GetString("userId")
	intr, err := i.interestSvc.GetInterest(userId)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "userId is not match with our resource")
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		errServerResp(c, err)
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
		Bio *string `json:"bio" binding:"required,max=300"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Bio": "at least an empty string and maximal character length is less than 300",
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
	}
	bio.Bio = strings.TrimSpace(*input.Bio)

	err = i.interestSvc.CreateNewBio(&bio)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefUserIdField) {
			errNotFoundResp(c, "userId is not match with our resource")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainUserId) {
			errUnprocessableEntityResp(c, "interest with this user id is already created")
			return
		}
		errServerResp(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "interest bio created",
	})

}
func (i *interest) putInterestBioHandler(c *gin.Context) {
	userId := c.GetString("userId")
	bio, err := i.interestSvc.GetBio(userId)
	if err != nil {
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "userId is not match with our resource")
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		errServerResp(c, err)
		return
	}
	var input struct {
		Bio *string `json:"bio" binding:"required,max=300"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Bio": "required, maximal character length is less than 300",
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
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "userId is not found in our db!")
			return
		}
		errServerResp(c, err)
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
		Hobbies []string `json:"hobbies" binding:"required,max=10,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Hobbies": "each hobbies must be unique, less than 10 and has more than 2 and less than 50 character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
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
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errNotFoundResp(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			errUnprocessableEntityResp(c, "every hobbies must be unique")
			return
		}
		if errors.Is(err, service.ErrCheckConstrainHobbie) {
			errUnprocessableEntityResp(c, "hobbies must less than 10")
			return
		}
		errServerResp(c, err)
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
		Hobbies []domain.Hobbie `json:"hobbies" binding:"required,max=10,unique=Hobbie"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Hobbies": "each hobbies must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new hobbies",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.PutHobbies(interestId, input.Hobbies)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errNotFoundResp(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			errUnprocessableEntityResp(c, "all of the hobbies must be unique")
			return
		}
		if errors.Is(err, service.ErrCheckConstrainHobbie) {
			errUnprocessableEntityResp(c, "hobbies must less than 10")
			return
		}
		errServerResp(c, err)
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
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Ids": "each ids must be unique and uuid character",
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
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "one of the id is not found")
			return
		}
		errServerResp(c, err)
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
		MovieSeries []string `json:"movieSeries" binding:"required,max=10,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"MovieSeries": "each movieSeries must be unique, less than 10 and has more than 2 and less than 50 character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
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
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errNotFoundResp(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			errUnprocessableEntityResp(c, "every moviesSeries must be unique")
			return
		}
		if errors.Is(err, service.ErrCheckConstrainMovieSeries) {
			errUnprocessableEntityResp(c, "movieSeries must less than 10")
			return
		}
		errServerResp(c, err)
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
		MovieSeries []domain.MovieSerie `json:"movieSeries" binding:"required,max=10,unique=MovieSerie"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"MovieSeries": "each movieSeries must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new movieSeries",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.PutMovieSeries(interestId, input.MovieSeries)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errNotFoundResp(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			errUnprocessableEntityResp(c, "all of the movieSeries must be unique")

			return
		}
		if errors.Is(err, service.ErrCheckConstrainMovieSeries) {
			errUnprocessableEntityResp(c, "movieSeries must less than 10")
			return
		}
		errServerResp(c, err)
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
	interestId := c.GetString("interestId")
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Ids": "each ids must be unique and uuid character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = i.interestSvc.DeleteMovieSeries(interestId, input.Ids)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "one of the id is not found")
			return
		}
		errServerResp(c, err)
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
		Travels []string `json:"travels" binding:"required,max=10,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Travels": "each travels must be unique, less than 10 and has more than 2 and less than 50 character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
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
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errNotFoundResp(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			errUnprocessableEntityResp(c, "every travels must be unique")
			return
		}
		if errors.Is(err, service.ErrCheckConstrainTraveling) {
			errUnprocessableEntityResp(c, "travels must less than 10")
			return
		}
		errServerResp(c, err)
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
		Travels []domain.Travel `json:"travels" binding:"required,max=10,unique=Travel"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Travels": "each travels must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new travel.",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.PutTraveling(interestId, input.Travels)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errNotFoundResp(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			errUnprocessableEntityResp(c, "all of the travels must be unique")
			return
		}
		if errors.Is(err, service.ErrCheckConstrainTraveling) {
			errUnprocessableEntityResp(c, "travels must less than 10")
			return
		}

		errServerResp(c, err)
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
	interestId := c.GetString("interestId")
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Ids": "each ids must be unique and uuid character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = i.interestSvc.DeleteTravels(interestId, input.Ids)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "one of the id is not found")
			return
		}
		errServerResp(c, err)
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
		Sports []string `json:"sports" binding:"required,max=10,unique,dive,min=2,max=50"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Sports": "each sports must be unique, less than 10 and has more than 2 and less than 50 character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
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
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errNotFoundResp(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			errUnprocessableEntityResp(c, "every sports must be unique")
			return
		}
		if errors.Is(err, service.ErrCheckConstrainSports) {
			errUnprocessableEntityResp(c, "sports must less than 10")
			return
		}
		errServerResp(c, err)
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
		Sports []domain.Sport `json:"sports" binding:"required,max=10,unique=Sport"`
	}
	interestId := c.GetString("interestId")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Sports": "each sports must be unique, less than 10 and has more than 2 and less than 50 character. Id must match or empty when its new sports.",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	err = i.interestSvc.PutSports(interestId, input.Sports)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, service.ErrRefInterestField) {
			errNotFoundResp(c, "interestId is not found")
			return
		}
		if errors.Is(err, service.ErrUniqueConstrainInterestId) {
			errUnprocessableEntityResp(c, "all of the sports must be unique")
			return
		}
		if errors.Is(err, service.ErrCheckConstrainSports) {
			errUnprocessableEntityResp(c, "sports must less than 10")
			return
		}
		errServerResp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"sports": input.Sports,
		},
	})
}
func (i *interest) deleteInterestSportHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Ids []string `json:"ids" binding:"required,unique,dive,uuid"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"Ids": "each ids must be unique and uuid character",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	err = i.interestSvc.DeleteSports(interestId, input.Ids)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "one of the id is not found")
			return
		}
		errServerResp(c, err)
		return
	}

	c.JSON(http.StatusNoContent, gin.H{
		"status":  "success",
		"message": "sports is successfully deleted",
	})
}
