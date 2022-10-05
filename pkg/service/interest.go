package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

var (
	ErrInterestIdNotFound          = errors.New("database: not found interest_id")
	ErrNotUniqueInterestIdResource = errors.New("database: violates unique constrain on interest id")
)

type Interest interface {
	GetInterest(userId string) (*domain.Interest, error)

	CreateNewBio(intr *domain.Bio) error
	GetBio(userId string) (*domain.Bio, error)
	PutBio(bio *domain.Bio) error

	CreateNewHobbies(interestId string, hobbies []domain.Hobbie) error
	PutHobbies(interestId string, hobbies []domain.Hobbie) error
	DeleteHobbies(ids []string) error

	CreateNewMovieSeries(interestId string, movieSeries []domain.MovieSerie) error
	PutMovieSeries(interestId string, movieSeries []domain.MovieSerie) error
	DeleteMovieSeries(ids []string) error

	CreateNewTraveling(interestId string, travels []domain.Travel) error
	PutTraveling(interestId string, travels []domain.Travel) error
	DeleteTravels(ids []string) error

	CreateNewSports(interestId string, sports []domain.Sport) error
	PutSports(interestId string, sports []domain.Sport) error
	DeleteSports(ids []string) error
}

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
	rows, err := i.interestRepo.UpdateInterestHobbies(interestId, hobbies)
	if err != nil {
		return i.parsingError(err)
	}
	if rows == 0 {
		panic("rows affected should not be zero")
	}
	return nil
}

func (i *interest) DeleteHobbies(ids []string) error {
	_, err := i.interestRepo.DeleteInterestHobbies(ids)
	if err != nil {
		return i.parsingError(err)
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

func (i *interest) DeleteMovieSeries(ids []string) error {
	_, err := i.interestRepo.DeleteInterestMovieSeries(ids)
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

func (i *interest) DeleteTravels(ids []string) error {
	_, err := i.interestRepo.DeleteInterestTraveling(ids)
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

func (i *interest) DeleteSports(ids []string) error {
	_, err := i.interestRepo.DeleteInterestSports(ids)
	if err != nil {
		return i.parsingError(err)
	}
	return nil
}

func (*interest) parsingError(err error) error {
	var pqErr *pq.Error
	switch {
	case errors.Is(err, context.Canceled):
		return domain.ErrTooLongAccesingDB
	case errors.Is(err, sql.ErrNoRows):
		return domain.ErrResourceNotFound
	case errors.As(err, &pqErr):
		if pqErr.Code == "23503" {
			if strings.Contains(pqErr.Constraint, "interest_id") {
				return ErrInterestIdNotFound
			}
			if strings.Contains(pqErr.Constraint, "user_id") {
				return ErrUserIdField
			}
			return err
		}
		if pqErr.Code == "23505" {
			if strings.Contains(pqErr.Constraint, "interest_id") {
				return ErrNotUniqueInterestIdResource
			}
			return err
		}
	}
	return err
}
