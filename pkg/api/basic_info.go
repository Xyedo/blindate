package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/service"
)

type basicInfoSvc interface {
	CreateBasicInfo(bInfo *domain.BasicInfo) error
	GetBasicInfoByUserId(id string) (*domain.BasicInfo, error)
	UpdateBasicInfo(bInfo *domain.BasicInfo) error
}

func NewBasicInfo(basicInfoService basicInfoSvc) basicinfo {
	return basicinfo{
		basicinfoService: basicInfoService,
	}
}

type basicinfo struct {
	basicinfoService basicInfoSvc
}

func (b basicinfo) postBasicInfoHandler(c *gin.Context) {
	userId := c.GetString("userId")
	var input struct {
		Gender           string  `json:"gender" binding:"required,max=25"`
		FromLoc          *string `json:"fromLoc" binding:"omitempty,max=25"`
		Height           *int    `json:"height" binding:"omitempty,min=0,max=300"`
		EducationLevel   *string `json:"educationLevel" binding:"omitempty,max=49"`
		Drinking         *string `json:"drinking" binding:"omitempty,max=49"`
		Smoking          *string `json:"smoking" binding:"omitempty,max=49"`
		RelationshipPref *string `json:"relationshipPref" binding:"omitempty,max=49"`
		LookingFor       string  `json:"lookingFor" binding:"required,max=25"`
		Zodiac           *string `json:"zodiac" binding:"omitempty,max=50"`
		Kids             *int    `json:"kids" binding:"omitempty,min=0,max=100"`
		Work             *string `json:"work" binding:"omitempty,max=50"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"gender":           "required and must have an gender enums",
			"fromLoc":          "maximal character is 100",
			"height":           "must have valid height in cm",
			"educationLevel":   "maximal character is 50",
			"drinking":         "maximal character is 50",
			"smoking":          "maximal character is 50",
			"relationshipPref": "maximal character is 50",
			"lookingFor":       "required and must have an gender enums",
			"zodiac":           "must have zodiac enums",
			"kids":             "minimum is 0 and maximal number is 100",
			"work":             "maximal character is 50",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
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
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "basic info created!",
	})
}

func (b basicinfo) getBasicInfoHandler(c *gin.Context) {
	userId := c.GetString("userId")
	basicInfo, err := b.basicinfoService.GetBasicInfoByUserId(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"basicInfo": basicInfo,
		},
	})
}

func (b basicinfo) patchBasicInfoHandler(c *gin.Context) {
	userId := c.GetString("userId")
	basicInfo, err := b.basicinfoService.GetBasicInfoByUserId(userId)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	var input struct {
		Gender           *string `json:"gender" binding:"omitempty,max=25"`
		FromLoc          *string `json:"fromLoc" binding:"omitempty,max=25"`
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
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"gender":           "maximal character is 25",
			"fromLoc":          "maximal character is 100",
			"height":           "must have valid height in cm",
			"educationLevel":   "maximal character is 50",
			"drinking":         "maximal character is 50",
			"smoking":          "maximal character is 50",
			"relationshipPref": "maximal character is 50",
			"lookingFor":       "maximal character is 25",
			"zodiac":           "maximal character is 50",
			"kids":             "minimum is 0 and maximal number is 100",
			"work":             "maximal character is 50",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
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
		basicInfo.Kids = input.Kids
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
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "basic info updated!",
	})

}
func referencesDbErr(err error) map[string]any {
	res := make(map[string]any)
	switch {
	case errors.Is(err, service.ErrRefUserIdField):
		res["status"] = "fail"
		res["message"] = "please refer to the documentation"
		res["errors"] = map[string]string{"user_id": "not found"}
	case errors.Is(err, service.ErrRefGenderField):
		res["status"] = "fail"
		res["message"] = "please refer to the documentation"
		res["errors"] = map[string]string{"gender": "not valid enums"}
	case errors.Is(err, service.ErrRefEducationLevelField):
		res["status"] = "fail"
		res["message"] = "please refer to the documentation"
		res["errors"] = map[string]string{"educationLevel": "not valid enums"}
	case errors.Is(err, service.ErrRefDrinkingField):
		res["status"] = "fail"
		res["message"] = "please refer to the documentation"
		res["errors"] = map[string]string{"drinking": "not valid enums"}
	case errors.Is(err, service.ErrRefSmokingField):
		res["status"] = "fail"
		res["message"] = "please refer to the documentation"
		res["errors"] = map[string]string{"smoking": "not valid enums"}
	case errors.Is(err, service.ErrRefRelationshipPrefField):
		res["status"] = "fail"
		res["message"] = "please refer to the documentation"
		res["errors"] = map[string]string{"relationshipPref": "not valid enums"}
	case errors.Is(err, service.ErrRefLookingForField):
		res["status"] = "fail"
		res["message"] = "please refer to the documentation"
		res["errors"] = map[string]string{"lookingFor": "not valid enums"}
	case errors.Is(err, service.ErrRefZodiacField):
		res["status"] = "fail"
		res["message"] = "please refer to the documentation"
		res["errors"] = map[string]string{"zodiac": "not valid enums"}
	default:
		res = nil
	}
	return res
}
