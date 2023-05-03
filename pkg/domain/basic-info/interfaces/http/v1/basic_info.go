package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	basicinfo "github.com/xyedo/blindate/pkg/domain/basic-info"
	basicInfoDTOs "github.com/xyedo/blindate/pkg/domain/basic-info/dtos"
	"github.com/xyedo/blindate/pkg/infrastructure"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func New(config infrastructure.Config, basicInfoUsecase basicinfo.Usecase) *basicInfoH {
	return &basicInfoH{
		config:      config,
		basicInfoUC: basicInfoUsecase,
	}
}

type basicInfoH struct {
	config      infrastructure.Config
	basicInfoUC basicinfo.Usecase
}

func (b *basicInfoH) postBasicInfoHandler(c *gin.Context) {
	var request postBasicInfoRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	err = request.mod().validate()
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	err = b.basicInfoUC.Create(basicInfoDTOs.CreateBasicInfo{
		UserId:           c.GetString(constant.KeyRequestUserId),
		Gender:           request.Gender,
		FromLoc:          request.FromLoc,
		Height:           request.Height,
		EducationLevel:   request.EducationLevel,
		Drinking:         request.Drinking,
		Smoking:          request.Smoking,
		RelationshipPref: request.RelationshipPref,
		LookingFor:       request.LookingFor,
		Zodiac:           request.Zodiac,
		Kids:             request.Kids,
		Work:             request.Work,
	})
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "basic info successfully created",
	})
}

func (b *basicInfoH) getBasicInfoHandler(c *gin.Context) {
	var url struct {
		UserId string `uri:"userId" binding:"required,uuid4"`
	}
	err := c.ShouldBindUri(&url)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "must have uuid in uri!",
		})
		return
	}

	basicInfo, err := b.basicInfoUC.GetById(c.GetString(constant.KeyRequestUserId), url.UserId)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"basic_info": getBasicInfoResponse{
				Gender:           basicInfo.Gender,
				FromLoc:          basicInfo.FromLoc,
				Height:           basicInfo.Height,
				EducationLevel:   basicInfo.EducationLevel,
				Drinking:         basicInfo.Drinking,
				Smoking:          basicInfo.Smoking,
				RelationshipPref: basicInfo.RelationshipPref,
				LookingFor:       basicInfo.LookingFor,
				Zodiac:           basicInfo.Zodiac,
				Kids:             basicInfo.Kids,
				Work:             basicInfo.Work,
			},
		},
	})
}

func (b *basicInfoH) patchBasicInfoHandler(c *gin.Context) {
	var request patchBasicInfoRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	err = request.mod().validate()
	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	err = b.basicInfoUC.Update(basicInfoDTOs.UpdateBasicInfo{
		UserId:           c.GetString(constant.KeyRequestUserId),
		Gender:           request.Gender,
		FromLoc:          request.FromLoc,
		Height:           request.Height,
		EducationLevel:   request.EducationLevel,
		RelationshipPref: request.EducationLevel,
		Drinking:         request.Drinking,
		Smoking:          request.Smoking,
		LookingFor:       request.LookingFor,
		Zodiac:           request.Zodiac,
		Kids:             request.Kids,
		Work:             request.Work,
	})

	if err != nil {
		httperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "basic-info updated",
	})
}
