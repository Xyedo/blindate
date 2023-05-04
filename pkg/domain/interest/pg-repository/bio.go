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
							Error:    err,
							ErrorMap: map[string][]string{"user_id": {"value is not found"}}})
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
					ErrorMap: map[string][]string{
						"user_id": {"value not found"},
					},
				},
			)
		}
		return err
	}

	return nil
}

// InsertMovieSeriesByInterestId implements interest.Repository
func (*interestConn) InsertMovieSeriesByInterestId(string, []interestEntities.MovieSerie) error {
	panic("unimplemented")
}

// UpdateMovieSeriesByInterestId implements interest.Repository
func (*interestConn) UpdateMovieSeriesByInterestId(string, []interestEntities.MovieSerie) error {
	panic("unimplemented")
}

// DeleteMovieSeriesByInterestId implements interest.Repository
func (*interestConn) DeleteMovieSeriesByInterestId(string, []string) error {
	panic("unimplemented")
}

// InsertSportByInterestId implements interest.Repository
func (*interestConn) InsertSportByInterestId(string, []interestEntities.Sport) error {
	panic("unimplemented")
}

// UpdateSportByInterestId implements interest.Repository
func (*interestConn) UpdateSportByInterestId(string, []interestEntities.Sport) error {
	panic("unimplemented")
}

// DeleteSportByInterestId implements interest.Repository
func (*interestConn) DeleteSportByInterestId(string, []string) error {
	panic("unimplemented")
}

// InsertTravelingByInterestId implements interest.Repository
func (*interestConn) InsertTravelingByInterestId(string, []interestEntities.Travel) error {
	panic("unimplemented")
}

// UpdateTravelingByInterestId implements interest.Repository
func (*interestConn) UpdateTravelingByInterestId(string, []interestEntities.Travel) error {
	panic("unimplemented")
}

// DeleteTravelingByInterestId implements interest.Repository
func (*interestConn) DeleteTravelingByInterestId(string, []string) error {
	panic("unimplemented")
}

// GetMovieSeriesByInterestId implements interest.Repository
func (*interestConn) GetMovieSeriesByInterestId(string) ([]interestEntities.MovieSerie, error) {
	panic("unimplemented")
}

// GetSportByInterestId implements interest.Repository
func (*interestConn) GetSportByInterestId(string) ([]interestEntities.Sport, error) {
	panic("unimplemented")
}

// GetTravelingByInterestId implements interest.Repository
func (*interestConn) GetTravelingByInterestId(string) ([]interestEntities.Travel, error) {
	panic("unimplemented")
}
