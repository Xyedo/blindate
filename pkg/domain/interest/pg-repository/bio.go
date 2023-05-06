package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	"github.com/xyedo/blindate/pkg/domain/interest"
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

func New(db *sqlx.DB) interest.Repository {
	return &interestConn{
		conn: db,
	}
}

type interestConn struct {
	conn *sqlx.DB
}

// InsertBio implements interest.Repository
func (i *interestConn) InsertBio(bio interestEntities.Bio) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnedBio string
	err := i.conn.GetContext(ctx, &returnedBio, insertBio, bio.UserId, bio.Bio, time.Now())
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return "", apperror.Timeout(apperror.Payload{Error: err})
		}

		if errors.Is(err, sql.ErrNoRows) {
			return "", apperror.NotFound(apperror.Payload{Error: err})
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				if strings.Contains(pqErr.Constraint, "user_id") {
					return "", apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string]string{
								"user_id": "value is not found",
							},
						})
				}
			}
			if pqErr.Code == "23505" {
				if strings.Contains(pqErr.Constraint, "interests_user_id_key") {
					return "", apperror.Conflicted(apperror.Payload{
						Error:   err,
						Message: "bio already inserted",
					})
				}
			}
			return "", err
		}
		return "", err
	}
	return returnedBio, nil
}

// GetBioByUserId implements interest.Repository
func (i *interestConn) GetBioByUserId(id string) (interestEntities.Bio, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var bio interestEntities.Bio
	err := i.conn.GetContext(ctx, &bio, getBioByUserId, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return interestEntities.Bio{}, apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return interestEntities.Bio{}, apperror.NotFound(apperror.Payload{Error: err, Message: "bio is not found"})
		}
		return interestEntities.Bio{}, err
	}

	return bio, nil
}

// UpdateBio implements interest.Repository
func (i *interestConn) UpdateBio(bio interestEntities.Bio) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnedId string
	err := i.conn.GetContext(
		ctx,
		&returnedId,
		updateInterestBio,
		bio.Bio,
		time.Now(),
		bio.UserId,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.UnprocessableEntity(
				apperror.PayloadMap{
					Error: err,
					ErrorMap: map[string]string{
						"user_id": "value not found",
					},
				},
			)
		}
		return err
	}

	return nil
}
