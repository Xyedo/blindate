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
	"github.com/xyedo/blindate/internal/tx"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	interestEntities "github.com/xyedo/blindate/pkg/domain/interest/entities"
)

var errInvalidBulkOperation = errors.New("invalid bulk operation")

// InsertHobbiesByInterestId implements interest.Repository
func (i *interestConn) InsertHobbiesByInterestId(id string, hobbies []interestEntities.Hobbie) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, len(hobbies))
	args = append(args, id)

	for i := range hobbies {
		param := i * 2
		stmt.WriteString(fmt.Sprintf("($%d, $%d, $%d),", 1, param+2, param+3))
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

		var returnedInterestId string
		err = tx.GetContext(
			ctx,
			&returnedInterestId,
			hobbiesStatistic,
			len(returnedIds),
			id,
		)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.UnprocessableEntity(
				apperror.PayloadMap{
					Error: err,
					ErrorMap: map[string][]string{
						"interest_id": {"value not found"},
					},
				})
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23503":
				if strings.Contains(pqErr.Constraint, "interest_id") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string][]string{
								"interest_id": {"value not found"},
							},
						})
				}
			case "23505":
				if strings.Contains(pqErr.Constraint, "hobbie_unique") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string][]string{
								"hobbie": {"every value must be unique"},
							},
						},
					)
				}
			case "23514":
				if strings.Contains(pqErr.Constraint, "hobbie_count") {
					return apperror.UnprocessableEntity(
						apperror.PayloadMap{
							Error: err,
							ErrorMap: map[string][]string{
								"hobbie": {"value must be less than 10"},
							},
						})
				}
			}
		}
		return err
	}

	return nil
}

// GetHobbiesByInterestId implements interest.Repository
func (i *interestConn) GetHobbiesByInterestId(id string) ([]interestEntities.Hobbie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var hobbies []interestEntities.Hobbie
	err := i.conn.GetContext(ctx, &hobbies, getHobbies, id)
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
func (i *interestConn) UpdateHobbiesByInterestId(id string, hobbies []interestEntities.Hobbie) error {
	stmt := new(strings.Builder)
	args := make([]any, 0, 2*len(hobbies))

	for i := range hobbies {
		param := i * 2
		stmt.WriteString(fmt.Sprintf("($%d, $%d),", param+1, param+2))
		args = append(args, hobbies[i].Id, hobbies[i].Hobbie)
	}
	statement := stmt.String()
	query := fmt.Sprintf(updateHobbies, statement[:len(statement)-1])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var returnedIds []string
	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		err := tx.GetContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(hobbies) {
			return errInvalidBulkOperation
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) {
			return apperror.BadPayload(apperror.Payload{Error: err, Message: "ids not found"})
		}
		if errors.Is(err, errInvalidBulkOperation) {
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
							ErrorMap: map[string][]string{
								"hobbies": {"every value must be unique"},
							},
						},
					)
				}
			}
			return err

		}
	}

	return nil
}

// DeleteHobbiesByInterestId implements interest.Repository
func (i *interestConn) DeleteHobbiesByInterestId(interestId string, hobbieIds []string) error {
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

	var returnedIds []string
	err := i.execTransaction(ctx, func(tx *sqlx.Tx) error {
		err := tx.SelectContext(ctx, &returnedIds, query, args...)
		if err != nil {
			return err
		}

		if len(returnedIds) != len(hobbieIds) {
			return errInvalidBulkOperation
		}

		var retInterestId string
		err = tx.GetContext(ctx, &retInterestId, hobbiesStatistic, -len(returnedIds), interestId)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return apperror.Timeout(apperror.Payload{Error: err})
		}
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, errInvalidBulkOperation) {
			return apperror.BadPayload(apperror.Payload{Error: err, Message: "ids is not valid"})
		}

		return err
	}

	return nil
}

func (i *interestConn) execTransaction(ctx context.Context, cb func(tx *sqlx.Tx) error) error {
	return tx.ExecGeneric(i.conn, ctx, cb, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
}
