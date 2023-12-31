package entities

import (
	"time"

	userentities "github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/pkg/optional"
)

type Match struct {
	Id                    string
	RequestFrom           string
	RequestTo             string
	RequestStatus         MatchStatus
	AcceptedAt            optional.Time
	RevealStatus          optional.Option[MatchStatus]
	RevealedDeclinedCount optional.Int32
	RevealedAt            optional.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
	UpdatedBy             optional.String
	Version               int
}

type MatchUser struct {
	Distance float64
	userentities.UserDetail
}
