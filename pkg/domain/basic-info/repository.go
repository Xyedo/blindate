package basicinfo

import (
	basicInfoEntities "github.com/xyedo/blindate/pkg/domain/basic-info/entities"
)

type Repository interface {
	InsertBasicInfo(basicInfoEntities.BasicInfo) error
	GetBasicInfoByUserId(string) (basicInfoEntities.BasicInfo, error)
	UpdateBasicInfo(bInfo basicInfoEntities.BasicInfo) error
}
