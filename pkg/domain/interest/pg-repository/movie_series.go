package repository

import (
	"context"
	"database/sql"
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var lengthInMovieSeriesDB int
	err := i.conn.GetContext(
		ctx,
		&lengthInMovieSeriesDB,
		checkInsertMovieSeriesValid,
		interestId,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		return err
	}

	if lengthInMovieSeriesDB+newMovieSeriesLength > 10 {
		return apperror.UnprocessableEntity(apperror.PayloadMap{
			ErrorMap: map[string]string{
				"movie_series": fmt.Sprintf(
					"new %d movie_series already surpassed the movie_series limit",
					newMovieSeriesLength),
			},
		})
	}

	return nil
}

// GetMovieSeriesByInterestId implements interest.Repository
func (i *interestConn) GetMovieSeriesByInterestId(id string) (
	[]interestEntities.MovieSerie,
	error,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var movieSeries []interestEntities.MovieSerie
	err := i.conn.SelectContext(ctx, &movieSeries, getMovieSeries, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return movieSeries, nil
}

// UpdateMovieSeriesByInterestId implements interest.Repository
func (i *interestConn) UpdateMovieSeriesByInterestId(
	id string,
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
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.BadPayload(apperror.Payload{
				Error:   err,
				Message: "ids not found",
			})
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
func (i *interestConn) DeleteMovieSeriesByInterestId(
	id string,
	movieIds []string,
) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, len(movieIds))

	for i, id := range movieIds {
		stmt.WriteString(fmt.Sprintf("$%d,", i+1))
		args = append(args, id)
	}
	statement := stmt.String()
	query := fmt.Sprintf(deleteMovieSeries, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(movieIds) {
			return transaction.ErrInvalidBulkOperation
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.UnprocessableEntity(apperror.PayloadMap{
				Error: err,
				ErrorMap: map[string]string{
					"movie_serie_ids": "all of the value not found",
				},
			})
		}
		if errors.Is(err, transaction.ErrInvalidBulkOperation) {
			return apperror.BadPayload(apperror.Payload{
				Error:   err,
				Message: "ids is not valid",
			})
		}

		return err
	}

	return nil
}
