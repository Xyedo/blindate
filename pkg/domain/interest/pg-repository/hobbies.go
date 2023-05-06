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

// InsertHobbiesByInterestId implements interest.Repository
func (i *interestConn) InsertHobbiesByInterestId(
	id string,
	hobbies []interestEntities.Hobbie,
) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, 2*len(hobbies))
	args = append(args, id)

	for i := range hobbies {
		param := i * 2
		stmt.WriteString(
			fmt.Sprintf("($%d, $%d, $%d),", 1, param+2, param+3),
		)
		newId := uuid.New()
		args = append(args, newId, hobbies[i].Hobbie)
		hobbies[i].Id = newId.String()
	}

	statement := stmt.String()
	query := fmt.Sprintf(insertHobbies, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(hobbies) {
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
				if strings.Contains(pqErr.Constraint, "hobbie_unique") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string]string{
								"hobbies": "every value must be unique",
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

// CheckInsertHobbiesValid implements interest.Repository
func (i *interestConn) CheckInsertHobbiesValid(
	interestId string,
	newHobbiesLength int,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var lengthInHobbiesDB int
	err := i.conn.GetContext(
		ctx,
		&lengthInHobbiesDB,
		checkInsertHobbiesValid,
		interestId,
	)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		return err
	}

	if lengthInHobbiesDB+newHobbiesLength > 10 {
		return apperror.UnprocessableEntity(apperror.PayloadMap{
			ErrorMap: map[string]string{
				"hobbies": fmt.Sprintf(
					"new %d hobbies already surpassed the hobbies limit",
					newHobbiesLength),
			},
		})
	}

	return nil
}

// GetHobbiesByInterestId implements interest.Repository
func (i *interestConn) GetHobbiesByInterestId(id string) ([]interestEntities.Hobbie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var hobbies []interestEntities.Hobbie
	err := i.conn.SelectContext(ctx, &hobbies, getHobbies, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return hobbies, nil
}

// UpdateHobbiesByInterestId implements interest.Repository
func (i *interestConn) UpdateHobbiesByInterestId(
	id string,
	hobbies []interestEntities.Hobbie,
) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, 2*len(hobbies))

	for i := range hobbies {
		param := i * 2
		stmt.WriteString(fmt.Sprintf("($%d::uuid, $%d),", param+1, param+2))
		args = append(args, hobbies[i].Id, hobbies[i].Hobbie)
	}
	statement := stmt.String()
	query := fmt.Sprintf(updateHobbies, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(hobbies) {
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
				Message: "ids not found"},
			)
		}
		if errors.Is(err, transaction.ErrInvalidBulkOperation) {
			return apperror.BadPayload(apperror.Payload{Error: err, Message: "one of the id is not found"})
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				if strings.Contains(pqErr.Constraint, "hobbie_unique") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string]string{
								"hobbies": "every value must be unique",
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

// DeleteHobbiesByInterestId implements interest.Repository
func (i *interestConn) DeleteHobbiesByInterestId(
	id string,
	hobbieIds []string,
) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, len(hobbieIds))

	for i, id := range hobbieIds {
		stmt.WriteString(fmt.Sprintf("$%d,", i+1))
		args = append(args, id)
	}

	statement := stmt.String()
	query := fmt.Sprintf(deleteHobbies, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		var returnedIds []string
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(hobbieIds) {
			return transaction.ErrInvalidBulkOperation
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) ||
			errors.Is(err, transaction.ErrInvalidBulkOperation) {
			return apperror.BadPayload(apperror.Payload{
				Error:   err,
				Message: "ids is not valid"})
		}

		return err
	}

	return nil
}

func (i *interestConn) execTransaction(ctx context.Context, cb func(tx *sqlx.Tx) error) error {
	return transaction.ExecGeneric(i.conn, ctx, cb, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
}
