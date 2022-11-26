package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

var (
	ErrRefInterestField          = fmt.Errorf("%w::interest_id", domain.ErrRefNotFound23503)
	ErrUniqueConstrainInterestId = fmt.Errorf("%w::interest_id", domain.ErrUniqueConstraint23505)
	ErrUniqueConstrainUserId     = fmt.Errorf("%w::user_id", domain.ErrUniqueConstraint23505)
	ErrCheckConstrainHobbie      = fmt.Errorf("database: over than 10 unit")
	ErrCheckConstrainMovieSeries = fmt.Errorf("database: over than 10 unit")
	ErrCheckConstrainTraveling   = fmt.Errorf("database: over than 10 unit")
	ErrCheckConstrainSports      = fmt.Errorf("database: over than 10 unit")
)

func NewInterest(intrRepo repository.Interest) *interest {
	return &interest{
		interestRepo: intrRepo,
	}
}

type interest struct {
	interestRepo repository.Interest
}

func (i *interest) GetInterest(userId string) (*domain.Interest, error) {
	intr, err := i.interestRepo.GetInterest(userId)
	if err != nil {
		err = i.parsingError(err)
		return nil, err
	}
	return intr, nil
}

func (i *interest) CreateNewBio(intr *domain.Bio) error {
	err := i.interestRepo.InsertInterestBio(intr)
	if err != nil {
		return i.parsingError(err)
	}
	err = i.interestRepo.InsertNewStats(intr.Id)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *interest) GetBio(userId string) (*domain.Bio, error) {
	bio, err := i.interestRepo.SelectInterestBio(userId)
	if err != nil {
		err = i.parsingError(err)
		return nil, err
	}
	return bio, nil
}

func (i *interest) PutBio(bio *domain.Bio) error {
	err := i.interestRepo.UpdateInterestBio(bio)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}

func (i *interest) CreateNewHobbies(interestId string, hobbies []domain.Hobbie) error {
	err := i.interestRepo.InsertInterestHobbies(interestId, hobbies)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}

func (i *interest) PutHobbies(interestId string, hobbies []domain.Hobbie) error {
	_, err := i.interestRepo.UpdateInterestHobbies(interestId, hobbies)
	if err != nil {
		return i.parsingError(err)
	}

	return nil
}

func (i *interest) DeleteHobbies(interestId string, ids []string) error {
	rows, err := i.interestRepo.DeleteInterestHobbies(interestId, ids)
	if err != nil {
		return i.parsingError(err)
	}
	if rows == 0 {
		return domain.ErrResourceNotFound
	}
	return nil
}

func (i *interest) CreateNewMovieSeries(interestId string, movieSeries []domain.MovieSerie) error {
	err := i.interestRepo.InsertInterestMovieSeries(interestId, movieSeries)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *interest) PutMovieSeries(interestId string, movieSeries []domain.MovieSerie) error {
	rows, err := i.interestRepo.UpdateInterestMovieSeries(interestId, movieSeries)
	if err != nil {
		return i.parsingError(err)
	}
	if rows == 0 {
		panic("rows affected should not be zero")
	}
	return nil
}

func (i *interest) DeleteMovieSeries(interestId string, ids []string) error {
	_, err := i.interestRepo.DeleteInterestMovieSeries(interestId, ids)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}

func (i *interest) CreateNewTraveling(interestId string, travels []domain.Travel) error {
	err := i.interestRepo.InsertInterestTraveling(interestId, travels)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *interest) PutTraveling(interestId string, travels []domain.Travel) error {
	rows, err := i.interestRepo.UpdateInterestTraveling(interestId, travels)
	if err != nil {
		return i.parsingError(err)
	}
	if rows == 0 {
		panic("rows affected should not be zero")
	}
	return nil
}

func (i *interest) DeleteTravels(interestId string, ids []string) error {
	_, err := i.interestRepo.DeleteInterestTraveling(interestId, ids)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *interest) CreateNewSports(interestId string, sports []domain.Sport) error {
	err := i.interestRepo.InsertInterestSports(interestId, sports)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}
func (i *interest) PutSports(interestId string, sports []domain.Sport) error {
	rows, err := i.interestRepo.UpdateInterestSport(interestId, sports)
	if err != nil {
		return i.parsingError(err)
	}
	if rows == 0 {
		panic("rows affected should not be zero")
	}
	return nil
}

func (i *interest) DeleteSports(interestId string, ids []string) error {
	_, err := i.interestRepo.DeleteInterestSports(interestId, ids)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}

func (*interest) parsingError(err error) error {
	var pqErr *pq.Error
	switch {
	case errors.Is(err, context.Canceled):
		return domain.ErrTooLongAccessingDB
	case errors.Is(err, sql.ErrNoRows):
		return domain.ErrResourceNotFound
	case errors.As(err, &pqErr):
		if pqErr.Code == "23503" {
			if strings.Contains(pqErr.Constraint, "interest_id") {
				return ErrRefInterestField
			}
			if strings.Contains(pqErr.Constraint, "user_id") {
				return ErrRefUserIdField
			}
			return err
		}
		if pqErr.Code == "23505" {
			if strings.Contains(pqErr.Constraint, "interest_id") {
				return ErrUniqueConstrainInterestId
			}
			if strings.Contains(pqErr.Constraint, "user_id") {
				return ErrUniqueConstrainUserId
			}
			return err
		}
		if pqErr.Code == "24514" {
			if strings.Contains(pqErr.Constraint, "hobbie_count") {
				return ErrCheckConstrainHobbie
			}
			if strings.Contains(pqErr.Constraint, "movie_serie_count") {
				return ErrCheckConstrainMovieSeries
			}
			if strings.Contains(pqErr.Constraint, "traveling_count") {
				return ErrCheckConstrainTraveling
			}
			if strings.Contains(pqErr.Constraint, "sport_count") {
				return ErrCheckConstrainSports
			}
			return err
		}
	}
	return err
}
