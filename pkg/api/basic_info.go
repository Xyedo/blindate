package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

type basicInfoSvc interface {
	CreateBasicInfo(bInfo domain.BasicInfo) error
	GetBasicInfoByUserId(id string) (domain.BasicInfo, error)
	UpdateBasicInfo(userId string, newBasicInfo domain.UpdateBasicInfo) error
}

func NewBasicInfo(basicInfoService basicInfoSvc) *BasicInfo {
	return &BasicInfo{
		basicinfoService: basicInfoService,
	}
}

type BasicInfo struct {
	basicinfoService basicInfoSvc
}

func (b *BasicInfo) postBasicInfoHandler(c *gin.Context) {
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

	err := b.basicinfoService.CreateBasicInfo(basicInfo)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "basic info created!",
	})
}

func (b *BasicInfo) getBasicInfoHandler(c *gin.Context) {
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

func (b *BasicInfo) patchBasicInfoHandler(c *gin.Context) {
	var inputBasicInfo domain.UpdateBasicInfo
	err := c.ShouldBindJSON(&inputBasicInfo)
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
	userId := c.GetString("userId")
	err = b.basicinfoService.UpdateBasicInfo(userId, inputBasicInfo)
	if err != nil {
		jsonHandleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "basic info updated!",
	})

}
