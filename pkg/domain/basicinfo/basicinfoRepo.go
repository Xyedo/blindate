package basicinfo

import basicInfoEntity "github.com/xyedo/blindate/pkg/domain/basicinfo/entities"

type Repository interface {
	InsertBasicInfo(basicinfo basicInfoEntity.DAO) error
	GetBasicInfoByUserId(id string) (basicInfoEntity.DAO, error)
	UpdateBasicInfo(bInfo basicInfoEntity.DAO) error
}
