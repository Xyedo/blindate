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

// InsertTravelingByInterestId implements interest.Repository
func (i *interestConn) InsertTravelingByInterestId(
	id string,
	travels []interestEntities.Travel,
) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, 2*len(travels))
	args = append(args, id)

	for i := range travels {
		param := i * 2
		stmt.WriteString(
			fmt.Sprintf("($%d, $%d, $%d),", 1, param+2, param+3),
		)
		newId := uuid.New()
		args = append(args, newId, travels[i].Travel)
		travels[i].Id = newId.String()
	}

	statement := stmt.String()
	query := fmt.Sprintf(insertTravels, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(travels) {
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
				if strings.Contains(pqErr.Constraint, "travel_unique") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string]string{
								"travels": "every value must be unique",
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

// CheckInsertTravelingValid implements interest.Repository
func (i *interestConn) CheckInsertTravelingValid(
	interestId string,
	newTravelsLength int,
) error {
	return i.checkInsertValueValid(
		interestId,
		"travels",
		checkInsertTravelsValid,
		newTravelsLength,
	)
}

// GetTravelingByInterestId implements interest.Repository
func (i *interestConn) GetTravelingByInterestId(id string) (
	[]interestEntities.Travel,
	error,
) {
	return getValuesByInterestId[interestEntities.Travel](
		i.conn,
		id,
		getTravels,
	)
}

// UpdateTravelingByInterestId implements interest.Repository
func (i *interestConn) UpdateTraveling(travels []interestEntities.Travel) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, 2*len(travels))

	for i := range travels {
		param := i * 2
		stmt.WriteString(fmt.Sprintf("($%d::uuid, $%d),", param+1, param+2))
		args = append(args, travels[i].Id, travels[i].Travel)
	}
	statement := stmt.String()
	query := fmt.Sprintf(updateTravels, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(travels) {
			return transaction.ErrInvalidBulkOperation
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, transaction.ErrInvalidBulkOperation) {
			return apperror.BadPayload(apperror.Payload{Error: err, Message: "one of the id is not found"})
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				if strings.Contains(pqErr.Constraint, "travel_unique") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string]string{
								"travels": "every value must be unique",
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

// DeleteTravelingByInterestId implements interest.Repository
func (i *interestConn) DeleteTravelingByIDs(travelIds []string) error {
	return i.deleteValuesByIDs(travelIds, deleteTravels, "travel_ids")
}
