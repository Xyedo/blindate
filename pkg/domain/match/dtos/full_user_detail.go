package matchDTOs

import (
	basicInfoDTOs "github.com/xyedo/blindate/pkg/domain/basic-info/dtos"
	interestDTOs "github.com/xyedo/blindate/pkg/domain/interest/dtos"
	userDTOs "github.com/xyedo/blindate/pkg/domain/user/dtos"
)

type FullUserDetail struct {
	userDTOs.User
	basicInfoDTOs.BasicInfo
	interestDTOs.InterestDetail
}
