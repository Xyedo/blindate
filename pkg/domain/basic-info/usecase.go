package basicinfo

import (
	basicInfoDTOs "github.com/xyedo/blindate/pkg/domain/basic-info/dtos"
	basicInfoEntities "github.com/xyedo/blindate/pkg/domain/basic-info/entities"
)

type Usecase interface {
	Create(basicInfoDTOs.CreateBasicInfo) error
	GetById(string) (basicInfoEntities.BasicInfo, error)
	Update(basicInfoDTOs.UpdateBasicInfo) error
}
