package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

var (
	ErrRefInterestField           = domain.WrapWithNewError(fmt.Errorf("%w::interest_id", domain.ErrRefNotFound23503), http.StatusNotFound, "interestId is not found")
	ErrUniqueConstrainInterestId  = domain.WrapWithNewError(fmt.Errorf("%w::interest_id", domain.ErrUniqueConstraint23505), http.StatusUnprocessableEntity, "interest with this user id is already created")
	ErrUniqueConstrainHobbies     = domain.WrapWithNewError(fmt.Errorf("%w::hobbies", domain.ErrUniqueConstraint23505), http.StatusUnprocessableEntity, "every hobbies must be unique")
	ErrUniqueConstrainMovieSeries = domain.WrapWithNewError(fmt.Errorf("%w::movieSeries", domain.ErrUniqueConstraint23505), http.StatusUnprocessableEntity, "every moviesSeries must be unique")
	ErrUniqueConstrainTraveling   = domain.WrapWithNewError(fmt.Errorf("%w::traveling", domain.ErrUniqueConstraint23505), http.StatusUnprocessableEntity, "every travels must be unique")
	ErrUniqueConstrainSport       = domain.WrapWithNewError(fmt.Errorf("%w::sport", domain.ErrUniqueConstraint23505), http.StatusUnprocessableEntity, "every sports must be unique")
	ErrCheckConstrainHobbie       = domain.WrapWithNewError(errors.New("database: over than 10 unit"), http.StatusUnprocessableEntity, "hobbies must less than 10")
	ErrCheckConstrainMovieSeries  = domain.WrapWithNewError(errors.New("database: over than 10 unit"), http.StatusUnprocessableEntity, "movieSeries must less than 10")
	ErrCheckConstrainTraveling    = domain.WrapWithNewError(errors.New("database: over than 10 unit"), http.StatusUnprocessableEntity, "travels must less than 10")
	ErrCheckConstrainSports       = domain.WrapWithNewError(errors.New("database: over than 10 unit"), http.StatusUnprocessableEntity, "sports must less than 10")
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
		return domain.WrapError(err, domain.ErrTooLongAccessingDB)
	case errors.Is(err, sql.ErrNoRows):
		return domain.WrapError(err, domain.ErrResourceNotFound)
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
			if strings.Contains(pqErr.Constraint, "interests_user_id_key") {
				return ErrUniqueConstrainInterestId
			}
			if strings.Contains(pqErr.Constraint, "hobbie_unique") {
				return ErrUniqueConstrainHobbies
			}
			if strings.Contains(pqErr.Constraint, "movie_serie_unique") {
				return ErrUniqueConstrainMovieSeries
			}
			if strings.Contains(pqErr.Constraint, "travel_unique") {
				return ErrUniqueConstrainTraveling
			}
			if strings.Contains(pqErr.Constraint, "sport_unique") {
				return ErrUniqueConstrainSport
			}
			if strings.Contains(pqErr.Constraint, "user_id") {
				return domain.ErrUniqueConstraint23505
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
