package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/xyedo/blindate/internal/transaction"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

// InsertMovieSeriesByInterestId implements interest.Repository
func (i *interestConn) InsertMovieSeriesByInterestId(
	id string,
	movieSeries []interestEntities.MovieSerie,
) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, 2*len(movieSeries))
	args = append(args, id)

	for i := range movieSeries {
		param := i * 2
		stmt.WriteString(
			fmt.Sprintf("($%d, $%d, $%d),", 1, param+2, param+3),
		)
		newId := uuid.New()
		args = append(args, newId, movieSeries[i].MovieSerie)
		movieSeries[i].Id = newId.String()
	}

	statement := stmt.String()
	query := fmt.Sprintf(insertMovieSeries, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(movieSeries) {
			return transaction.ErrInvalidBulkOperation
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, transaction.ErrInvalidBulkOperation) {
			return apperror.BadPayload(apperror.Payload{
				Error:   err,
				Message: "invalid bulk operation",
			})
		}

		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23503":
				if strings.Contains(pqErr.Constraint, "interest_id") {
					return apperror.NotFound(
						apperror.Payload{
							Error:   err,
							Message: "interest_id not found",
						})
				}
			case "23505":
				if strings.Contains(pqErr.Constraint, "movie_serie_unique") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string]string{
								"movie_series": "every value must be unique",
							},
						},
					)
				}
			}
		}
		return err
	}

	return nil
}

// CheckInsertMovieSeriesValid implements interest.Repository
func (i *interestConn) CheckInsertMovieSeriesValid(
	interestId string,
	newMovieSeriesLength int,
) error {
	return i.checkInsertValueValid(
		interestId,
		"movie_series",
		checkInsertMovieSeriesValid,
		newMovieSeriesLength,
	)
}

// GetMovieSeriesByInterestId implements interest.Repository
func (i *interestConn) GetMovieSeriesByInterestId(id string) (
	[]interestEntities.MovieSerie,
	error,
) {
	return getValuesByInterestId[interestEntities.MovieSerie](
		i.conn,
		id,
		getMovieSeries,
	)
}

// UpdateMovieSeriesByInterestId implements interest.Repository
func (i *interestConn) UpdateMovieSeries(
	movieSeries []interestEntities.MovieSerie,
) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, 2*len(movieSeries))

	for i := range movieSeries {
		param := i * 2
		stmt.WriteString(fmt.Sprintf("($%d::uuid, $%d),", param+1, param+2))
		args = append(args, movieSeries[i].Id, movieSeries[i].MovieSerie)
	}
	statement := stmt.String()
	query := fmt.Sprintf(updateMovieSeries, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(movieSeries) {
			return transaction.ErrInvalidBulkOperation
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, transaction.ErrInvalidBulkOperation) {
			return apperror.BadPayload(apperror.Payload{
				Error:   err,
				Message: "one of the id is not found",
			})
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				if strings.Contains(pqErr.Constraint, "movie_serie_unique") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string]string{
								"movie_series": "every value must be unique",
							},
						},
					)
				}
			}
		}
		return err
	}

	return nil
}

// DeleteMovieSeriesByInterestId implements interest.Repository
func (i *interestConn) DeleteMovieSeriesByIDs(movieIds []string) error {
	return i.deleteValuesByIDs(
		movieIds,
		deleteMovieSeries,
		"movie_serie_ids",
	)
}
