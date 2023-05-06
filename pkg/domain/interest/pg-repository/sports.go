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

// InsertSportByInterestId implements interest.Repository
func (i *interestConn) InsertSportByInterestId(
	id string,
	sports []interestEntities.Sport,
) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, 2*len(sports))
	args = append(args, id)

	for i := range sports {
		param := i * 2
		stmt.WriteString(
			fmt.Sprintf("($%d, $%d, $%d),", 1, param+2, param+3),
		)
		newId := uuid.New()
		args = append(args, newId, sports[i].Sport)
		sports[i].Id = newId.String()
	}

	statement := stmt.String()
	query := fmt.Sprintf(insertSports, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(sports) {
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
				if strings.Contains(pqErr.Constraint, "sport_unique") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string]string{
								"sports": "every value must be unique",
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

// CheckInsertSportValid implements interest.Repository
func (i *interestConn) CheckInsertSportValid(
	interestId string,
	newSportsLength int,
) error {
	return i.checkInsertValueValid(
		interestId,
		"sports",
		checkInsertSportsValid,
		newSportsLength,
	)
}

// GetSportByInterestId implements interest.Repository
func (i *interestConn) GetSportByInterestId(id string) (
	[]interestEntities.Sport,
	error,
) {
	return getValuesByInterestId[interestEntities.Sport](
		i.conn,
		id,
		getSports,
	)
}

// UpdateSportByInterestId implements interest.Repository
func (i *interestConn) UpdateSport(
	sports []interestEntities.Sport,
) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, 2*len(sports))

	for i := range sports {
		param := i * 2
		stmt.WriteString(fmt.Sprintf("($%d::uuid, $%d),", param+1, param+2))
		args = append(args, sports[i].Id, sports[i].Sport)
	}
	statement := stmt.String()
	query := fmt.Sprintf(updateSports, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(sports) {
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
				if strings.Contains(pqErr.Constraint, "sport_unique") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string]string{
								"sports": "every value must be unique",
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

// DeleteSportByInterestId implements interest.Repository
func (i *interestConn) DeleteSportByIDs(sportIds []string) error {
	return i.deleteValuesByIDs(sportIds, deleteSports, "sport_ids")
}
