package basicinfo

import basicInfoEntity "github.com/xyedo/blindate/pkg/domain/basicinfo/entities"

type Repository interface {
	InsertBasicInfo(basicinfo basicInfoEntity.Dao) error
	GetBasicInfoByUserId(id string) (basicInfoEntity.Dao, error)
	UpdateBasicInfo(bInfo basicInfoEntity.Dao) error
}
