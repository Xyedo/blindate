package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	basicInfoEntity "github.com/xyedo/blindate/pkg/domain/basicinfo/entities"
)

type basicInfoSvc interface {
	CreateBasicInfo(bInfo basicInfoEntity.DTO) error
	GetBasicInfoByUserId(id string) (basicInfoEntity.DTO, error)
	UpdateBasicInfo(userId string, newBasicInfo basicInfoEntity.Update) error
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
	userId := c.GetString(keyUserId)
	var input basicInfoEntity.DTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"gender":           "required and the value must one of the gender enums",
			"fromLoc":          "maximal character is 100",
			"height":           "must have valid height in cm",
			"educationLevel":   "must one of the educationLevel enums",
			"drinking":         "must one of the drinking enums",
			"smoking":          "must one of the smoking enums",
			"relationshipPref": "must one of the relationshipPref enums",
			"lookingFor":       "required and the value must one of the lookingFor enums",
			"zodiac":           "must one of the zodiac enums",
			"kids":             "minimum is 0 and maximal number is 30",
			"work":             "maximal character is 50",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}

	basicInfo := basicInfoEntity.DTO{
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
	userId := c.GetString(keyUserId)
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
	var inputBasicInfo basicInfoEntity.Update
	err := c.ShouldBindJSON(&inputBasicInfo)
	if err != nil {
		errjson := jsonBindingErrResp(err, c, map[string]string{
			"gender":           "must one of the gender enums",
			"fromLoc":          "maximal character is 100",
			"height":           "must have valid height in cm",
			"educationLevel":   "must one of the educationLevel enums",
			"drinking":         "must one of the drinking enums",
			"smoking":          "must one of the smoking enums",
			"relationshipPref": "must one of the relationshipPref enums",
			"lookingFor":       "must one of the lookingFor enums",
			"zodiac":           "must one of the zodiac enums",
			"kids":             "minimum is 0 and maximal number is 30",
			"work":             "maximal character is 50",
		})
		if errjson != nil {
			errServerResp(c, err)
			return
		}
		return
	}
	userId := c.GetString(keyUserId)
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
