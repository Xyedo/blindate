package basicinfo

import (
	basicInfoDTOs "github.com/xyedo/blindate/pkg/domain/basic-info/dtos"
)

type Usecase interface {
	Create(basicInfoDTOs.CreateBasicInfo) error
	GetById(string, string) (basicInfoDTOs.BasicInfo, error)
	Update(basicInfoDTOs.UpdateBasicInfo) error
}
