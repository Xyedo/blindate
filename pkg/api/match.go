package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xyedo/blindate/pkg/domain"
)

type matchSvc interface {
	FindNewMatch(userId string) ([]domain.User, error)
}

func NewMatch(matchSvc matchSvc, interestSvc interestSvc, basicInfoSvc basicInfoSvc) *match {
	return &match{
		matchSvc:     matchSvc,
		interestSvc:  interestSvc,
		basicInfoSvc: basicInfoSvc,
	}
}

type match struct {
	matchSvc     matchSvc
	interestSvc  interestSvc
	basicInfoSvc basicInfoSvc
}

func (m *match) getNewMatchHandler(c *gin.Context) {
	result := []struct {
		domain.User
		domain.Interest
		domain.BasicInfo
	}{}
	userId := c.GetString("userId")
	matchUsers, err := m.matchSvc.FindNewMatch(userId)
	if err != nil {
		//TODO: better error handling
		errServerResp(c, err)
		return
	}
	for _, matchUser := range matchUsers {
		basicInfo, err := m.basicInfoSvc.GetBasicInfoByUserId(matchUser.ID)
		if err != nil {
			//TODO: better error handling
			errServerResp(c, err)
			return
		}

		interest, err := m.interestSvc.GetInterest(matchUser.ID)
		if err != nil {
			//TODO: better error handling
			errServerResp(c, err)
			return
		}
		result = append(result, struct {
			domain.User
			domain.Interest
			domain.BasicInfo
		}{
			User:      matchUser,
			Interest:  *interest,
			BasicInfo: *basicInfo,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"results": result,
		},
	})

}
