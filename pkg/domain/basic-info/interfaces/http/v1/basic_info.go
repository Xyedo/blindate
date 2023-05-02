package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/common/constant"
	basicinfo "github.com/xyedo/blindate/pkg/domain/basic-info"
	basicInfoDTOs "github.com/xyedo/blindate/pkg/domain/basic-info/dtos"
	httperror "github.com/xyedo/blindate/pkg/infrastructure/http/error"
)

func New(basicInfoUsecase basicinfo.Usecase) *basicInfoH {
	return &basicInfoH{
		basicInfoUC: basicInfoUsecase,
	}
}

type basicInfoH struct {
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

	userId := c.GetString(constant.KeyUserId)
	err = b.basicInfoUC.Create(basicInfoDTOs.CreateBasicInfo{
		UserId:           userId,
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
	userId := c.GetString(constant.KeyUserId)

	basicInfo, err := b.basicInfoUC.GetById(userId)
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
