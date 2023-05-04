package interestDTOs

import (
	"github.com/xyedo/blindate/internal/optional"
)

type UpdateBio struct {
	Id     string
	UserId string
	Bio    optional.String
}
