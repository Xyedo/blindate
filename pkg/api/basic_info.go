package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
	"github.com/xyedo/blindate/pkg/util"
)

func NewBasicInfo(basicInfoService service.BasicInfo) *basicinfo {
	return &basicinfo{
		basicinfoService: basicInfoService,
	}
}

type basicinfo struct {
	basicinfoService service.BasicInfo
}

func (b *basicinfo) postBasicInfoHandler(c *gin.Context) {
	userId := c.GetString("userId")
	var input struct {
		Gender           string  `json:"gender" binding:"required,max=25"`
		FromLoc          *string `json:"fromLoc" binding:"binding:required,max=25"`
		Height           *int    `json:"height" binding:"omitempty,min=0,max=300"`
		EducationLevel   *string `json:"educationLevel" binding:"omitempty,max=49"`
		Drinking         *string `json:"drinking" binding:"omitempty,max=49"`
		Smoking          *string `json:"smoking" binding:"omitempty,max=49"`
		RelationshipPref *string `json:"relationshipPref" binding:"omitempty,max=49"`
		LookingFor       string  `json:"lookingFor" binding:"required,max=25"`
		Zodiac           *string `json:"zodiac" binding:"omitempty,max=50"`
		Kids             int     `json:"kids" binding:"required,max=100"`
		Work             *string `json:"work" binding:"omitempty,max=50"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Gender":           "required and must have an gender enums",
			"FromLoc":          "maximal character is 100",
			"Height":           "must have valid height in cm",
			"EducationLevel":   "maximal character is 50",
			"Drinking":         "maximal character is 50",
			"Smoking":          "maximal character is 50",
			"RelationshipPref": "maximal character is 50",
			"LookingFor":       "required and must have an gender enums",
			"Zodiac":           "must have zodiac enums",
			"Kids":             "minimum is 0 and maximal number is 100",
			"Work":             "maximal character is 50",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}

	basicInfo := domain.BasicInfo{
		UserId:           userId,
		Gender:           input.Gender,
		FromLoc:          input.FromLoc,
		Height:           input.Height,
		EducationLevel:   input.EducationLevel,
		Drinking:         input.Drinking,
		Smoking:          input.Smoking,
		RelationshipPref: input.RelationshipPref,
		LookingFor:       input.LookingFor,
		Zodiac:           input.Zodiac,
		Kids:             input.Kids,
		Work:             input.Work,
	}

	err := b.basicinfoService.CreateBasicInfo(&basicInfo)
	if err != nil {
		res := referencesDbErr(err)
		if res != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, res)
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorRequestTimeout(c)
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "basic info created!",
	})
}

func (b *basicinfo) getBasicInfoHandler(c *gin.Context) {
	userId := c.GetString("userId")
	basicInfo, err := b.basicinfoService.GetBasicInfoByUserId(userId)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorRequestTimeout(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "user Id not match with our basic info")
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "succes",
		"data": gin.H{
			"basicInfo": basicInfo,
		},
	})
}

func (b *basicinfo) patchBasicInfoHandler(c *gin.Context) {
	userId := c.GetString("userId")
	basicInfo, err := b.basicinfoService.GetBasicInfoByUserId(userId)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorRequestTimeout(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errorResourceNotFound(c, "user Id not match with our basic info")
			return
		}
		errorServerResponse(c, err)
		return
	}
	var input struct {
		Gender           *string `json:"gender" binding:"omitempty,max=25"`
		FromLoc          *string `json:"fromLoc" binding:"binding:omitempty,max=25"`
		Height           *int    `json:"height" binding:"omitempty,min=0,max=300"`
		EducationLevel   *string `json:"educationLevel" binding:"omitempty,max=49"`
		Drinking         *string `json:"drinking" binding:"omitempty,max=49"`
		Smoking          *string `json:"smoking" binding:"omitempty,max=49"`
		RelationshipPref *string `json:"relationshipPref" binding:"omitempty,max=49"`
		LookingFor       *string `json:"lookingFor"  binding:"omitempty,max=25"`
		Zodiac           *string `json:"zodiac" binding:"omitempty,max=50"`
		Kids             *int    `json:"kids" binding:"omitempty,max=100"`
		Work             *string `json:"work" binding:"omitempty,max=50"`
	}
	err = c.ShouldBindJSON(&input)
	if err != nil {
		err1 := util.ReadJSONDecoderErr(err)
		if err1 != nil {
			errorJSONBindingResponse(c, err1)
			return
		}
		errMap := util.ReadValidationErr(err, map[string]string{
			"Gender":           "maximal character is 25",
			"FromLoc":          "maximal character is 100",
			"Height":           "must have valid height in cm",
			"EducationLevel":   "maximal character is 50",
			"Drinking":         "maximal character is 50",
			"Smoking":          "maximal character is 50",
			"RelationshipPref": "maximal character is 50",
			"LookingFor":       "maximal character is 25",
			"Zodiac":           "maximal character is 50",
			"Kids":             "minimum is 0 and maximal number is 100",
			"Work":             "maximal character is 50",
		})
		if errMap != nil {
			errorValidationResponse(c, errMap)
			return
		}
		errorServerResponse(c, err)
		return
	}
	if input.Gender != nil {
		basicInfo.Gender = *input.Gender
	}
	if input.FromLoc != nil {
		basicInfo.FromLoc = input.FromLoc
	}
	if input.Height != nil {
		basicInfo.Height = input.Height
	}
	if input.EducationLevel != nil {
		basicInfo.EducationLevel = input.EducationLevel
	}
	if input.Drinking != nil {
		basicInfo.Drinking = input.Drinking
	}
	if input.Smoking != nil {
		basicInfo.Smoking = input.Smoking
	}
	if input.RelationshipPref != nil {
		basicInfo.RelationshipPref = input.RelationshipPref
	}
	if input.LookingFor != nil {
		basicInfo.LookingFor = *input.LookingFor
	}
	if input.Zodiac != nil {
		basicInfo.Zodiac = input.Zodiac
	}
	if input.Kids != nil {
		basicInfo.Kids = *input.Kids
	}
	if input.Work != nil {
		basicInfo.Work = input.Work
	}
	err = b.basicinfoService.UpdateBasicInfo(basicInfo)
	if err != nil {
		res := referencesDbErr(err)
		if res != nil {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, res)
			return
		}
		if errors.Is(err, domain.ErrTooLongAccesingDB) {
			errorRequestTimeout(c)
			return
		}
		errorServerResponse(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "basic info updated!",
	})

}
func referencesDbErr(err error) map[string]string {
	res := map[string]string{}
	switch {
	case errors.Is(err, service.ErrUserIdField):
		res["status"] = "fail"
		res["message"] = "userid not found!"
	case errors.Is(err, service.ErrGenderField):
		res["status"] = "fail"
		res["message"] = "gender property is not valid enums, refer to the documentation!"
	case errors.Is(err, service.ErrEducationLevelField):
		res["status"] = "fail"
		res["message"] = "educationLevel property is not valid enums, refer to the documentation!"
	case errors.Is(err, service.ErrDrinkingField):
		res["status"] = "fail"
		res["message"] = "drinking property is not valid enums, refer to the documentation!"
	case errors.Is(err, service.ErrSmokingField):
		res["status"] = "fail"
		res["message"] = "smoking property is not valid enums, refer to the documentation!"
	case errors.Is(err, service.ErrrRelationshipPrefField):
		res["status"] = "fail"
		res["message"] = "relationshipPref property is not valid enums, refer to the documentation!"
	case errors.Is(err, service.ErrLookingForField):
		res["status"] = "fail"
		res["message"] = "lookingFor property is not valid enums, refer to the documentation!"
	case errors.Is(err, service.ErrZodiacField):
		res["status"] = "fail"
		res["message"] = "zodiac property is not valid enums, refer to the documentation!"
	default:
		res = nil
	}
	return res
}
