package repository

import (
	"context"

	"github.com/xyedo/blindate/internal/domain/user/entities"
	"github.com/xyedo/blindate/internal/infrastructure/pg"
)

func FindNonMatchClosestUser(ctx context.Context, conn pg.Querier, payload entities.FindClosestUser) ([]string, error) {
	const findClosestUserById = `
	SELECT
		ad.account_id,
		ad.geog <-> ST_GeomFromText($2) as distance
	FROM account_detail ad
	WHERE 
		ad.account_id != $1 AND
		NOT EXISTS (
			SELECT 1
			FROM match m
			WHERE 
				m.request_to = ad.account_id OR
				m.request_from = ad.account_id
		)
	ORDER BY distance
	OFFSET $3
	LIMIT $4
	`

	rows, err := conn.Query(ctx, findClosestUserById,
		payload.UserId,
		payload.Geog,
		payload.Pagination.Offset(),
		payload.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	closestUsers := make([]string, 0)
	for rows.Next() {
		var closestUser entities.ClosestUser
		err = rows.Scan(
			&closestUser.UserId,
			&closestUser.Distance,
		)
		if err != nil {
			return nil, err
		}

		closestUsers = append(closestUsers, closestUser.UserId)
	}

	return closestUsers, nil
}
