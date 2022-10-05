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

// 		MovieSeries []string `json:"movieSeries" binding:"omitempty,unique,dive,min=2,max=50"`
// 		Travels     []string `json:"travels" binding:"omitempty,unique,dive,min=2,max=50"`
// 		Sports      []string `json:"sports" binding:"omitempty,unique,dive,min=2,max=50"`

// "MovieSeries": "each movieSeries must be unique and has more than 2 and less than 50 character",
// "Travels":     "each travels must be unique and has more than 2 and less than 50 character",
// "Sports":      "each travels must be unique and has more than 2 and less than 50 character",
func (i *interest) GetInterestHandler(c *gin.Context) {
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

func (i *interest) PostInterestBioHandler(c *gin.Context) {
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
	// for _, hobbie := range input.Hobbies {
	// 	intr.Hobbies = append(intr.Hobbies, domain.Hobbie{
	// 		Hobbie: hobbie,
	// 	})
	// }
	// for _, movieSerie := range input.MovieSeries {
	// 	intr.MovieSeries = append(intr.MovieSeries, domain.MovieSerie{
	// 		MovieSerie: movieSerie,
	// 	})
	// }
	// for _, travel := range input.Travels {
	// 	intr.Travels = append(intr.Travels, domain.Travel{
	// 		Travel: travel,
	// 	})
	// }
	// for _, sport := range input.Sports {
	// 	intr.Sports = append(intr.Sports, domain.Sport{
	// 		Sport: sport,
	// 	})
	// }

	err = i.interestSvc.CreateNewBio(&bio)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorDeadLockResponse(c)
			return
		}
		if errors.Is(err, service.ErrUserIdField) {
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
func (i *interest) PutInterestBioHandler(c *gin.Context) {
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
		// Hobbies     []domain.Hobbie `json:"hobbies" binding:"omitempty,unique=Id,unique=Hobbie"`
		// MovieSeries []domain.MovieSerie `json:"movieSeries" binding:"omitempty,unique=Id,unique=MovieSerie"`
		// Travels     []domain.Travel `json:"travels" binding:"omitempty,unique=Id,unique=Travel"`
		// Sports      []domain.Sport `json:"sports" binding:"omitempty,unique=Id,unique=Sport"`
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
			// "Hobbies":     "each hobbies must be unique and has more than 2 and less than 50 character",
			// "MovieSeries": "each movieSeries must be unique and has more than 2 and less than 50 character",
			// "Travels":     "each travels must be unique and has more than 2 and less than 50 character",
			// "Sports":      "each travels must be unique and has more than 2 and less than 50 character",
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

func (i *interest) PostInterestHobbiesHandler(c *gin.Context) {
	interestId := c.GetString("interestId")
	var input struct {
		Hobbies []string `json:"hobbies" binding:"required,max=5,unique,dive,min=2,max=50"`
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
		if errors.Is(err, service.ErrInterestIdNotFound) {
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
