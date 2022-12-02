package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
	"github.com/xyedo/blindate/pkg/repository"
)

func NewMatch(matchRepo repository.Match, locationRepo repository.Location) *match {
	return &match{
		matchRepo:    matchRepo,
		locationRepo: locationRepo,
	}
}

type match struct {
	matchRepo    repository.Match
	locationRepo repository.Location
}

func (m *match) FindUserToMatch(userId string) ([]domain.BigUser, error) {
	userLoc, err := m.locationRepo.GetLocationByUserId(userId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccessingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		return nil, err
	}
	toUsers, err := m.locationRepo.GetClosestUser(userId, userLoc.Geog, 3)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccessingDB
		}
		return nil, err
	}

	return toUsers, nil
}
func (m *match) PostNewMatch(fromUserId, toUserId string, matchStatus domain.MatchStatus) (string, error) {
	id, err := m.matchRepo.InsertNewMatch(fromUserId, toUserId, matchStatus)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.Is(err, context.Canceled):
			return "", domain.ErrTooLongAccessingDB
		case errors.As(err, &pqErr):
			switch {
			case pqErr.Code == "23503":
				return "", domain.ErrRefNotFound23503
			case pqErr.Code == "23505":
				return "", domain.ErrUniqueConstraint23505
			default:
				return "", pqErr
			}
		default:
			return "", err
		}
	}
	return id, nil

}
func (m *match) GetMatchReqToUserId(userId string) ([]domain.MatchUser, error) {
	matcheds, err := m.matchRepo.SelectMatchReqToUserId(userId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccessingDB
		}
		return nil, err
	}

	return matcheds, nil
}

func (m *match) RequestChange(matchId string, matchStatus domain.MatchStatus) error {
	matchEntity, err := m.getMatchById(matchId)
	if err != nil {
		return err
	}
	switch matchStatus {
	case domain.Requested:
		if matchEntity.RequestStatus != string(domain.Unknown) {
			return ErrInvalidMatchStatus
		}
		matchEntity.RequestStatus = string(domain.Requested)
	case domain.Declined:
		matchEntity.RequestStatus = string(domain.Declined)
	case domain.Accepted:
		if matchEntity.RequestStatus != string(domain.Requested) {
			return ErrInvalidMatchStatus
		}
		matchEntity.RequestStatus = string(domain.Accepted)
	}
	err = m.updateMatch(*matchEntity)
	if err != nil {
		return err
	}
	return nil
}

func (m *match) RevealChange(matchId string, matchStatus domain.MatchStatus) error {
	matchEntity, err := m.getMatchById(matchId)
	if err != nil {
		return err
	}
	if matchEntity.RequestStatus != string(domain.Accepted) {
		return ErrInvalidMatchStatus
	}
	switch matchStatus {
	case domain.Requested:
		if matchEntity.RevealStatus != string(domain.Unknown) {
			return ErrInvalidMatchStatus
		}
		matchEntity.RevealStatus = string(domain.Requested)
	case domain.Declined:
		if matchEntity.RevealStatus == string(domain.Unknown) {
			return ErrInvalidMatchStatus
		}
		matchEntity.RevealStatus = string(domain.Declined)
	case domain.Accepted:
		if matchEntity.RevealStatus != string(domain.Requested) {
			return ErrInvalidMatchStatus
		}
		matchEntity.RevealStatus = string(domain.Accepted)
	}

	err = m.updateMatch(*matchEntity)
	if err != nil {
		return err
	}
	return nil
}

func (m *match) updateMatch(matchEntity entity.Match) error {
	err := m.matchRepo.UpdateMatchById(matchEntity)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.ErrResourceNotFound
		case errors.Is(err, context.Canceled):
			return domain.ErrTooLongAccessingDB
		default:
			return err
		}
	}
	return nil
}

func (m *match) getMatchById(matchId string) (*entity.Match, error) {
	matchEntity, err := m.matchRepo.GetMatchById(matchId)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrResourceNotFound
		case errors.Is(err, context.Canceled):
			return nil, domain.ErrTooLongAccessingDB
		default:
			return nil, err
		}
	}
	return matchEntity, nil
}
