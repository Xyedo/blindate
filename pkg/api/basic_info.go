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

func NewBasicInfo(basicInfoService basicInfoSvc) *basicinfo {
	return &basicinfo{
		basicinfoService: basicInfoService,
	}
}

type basicinfo struct {
	basicinfoService basicInfoSvc
}

func (b *basicinfo) postBasicInfoHandler(c *gin.Context) {
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
		if errors.Is(err, domain.ErrTooLongAccessingDB) {
			errResourceConflictResp(c)
			return
		}
		errServerResp(c, err)
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
		if errors.Is(err, domain.ErrTooLongAccessingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "user Id not match with our basic info")
			return
		}
		errServerResp(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"basicInfo": basicInfo,
		},
	})
}

func (b *basicinfo) patchBasicInfoHandler(c *gin.Context) {
	userId := c.GetString("userId")
	basicInfo, err := b.basicinfoService.GetBasicInfoByUserId(userId)
	if err != nil {
		if errors.Is(err, domain.ErrTooLongAccessingDB) {
			errResourceConflictResp(c)
			return
		}
		if errors.Is(err, domain.ErrResourceNotFound) {
			errNotFoundResp(c, "user Id not match with our basic info")
			return
		}
		errServerResp(c, err)
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
		if errors.Is(err, domain.ErrTooLongAccessingDB) {
			errResourceConflictResp(c)
			return
		}
		errServerResp(c, err)
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
	case errors.Is(err, service.ErrUniqueConstrainUserId):
		res["status"] = "fail"
		res["message"] = "basic_info with this user id is already created"
	default:
		res = nil
	}
	return res
}
