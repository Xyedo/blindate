package service

import (
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/domain/entity"
	"github.com/xyedo/blindate/pkg/repository"
)

func NewMatch(matchRepo repository.Match, locationRepo repository.Location) *Match {
	return &Match{
		matchRepo:    matchRepo,
		locationRepo: locationRepo,
	}
}

type Match struct {
	matchRepo    repository.Match
	locationRepo repository.Location
}

func (m *Match) FindUserToMatch(userId string) ([]domain.BigUser, error) {
	userLoc, err := m.locationRepo.GetLocationByUserId(userId)
	if err != nil {
		return nil, err
	}
	toUsers, err := m.locationRepo.GetClosestUser(userId, userLoc.Geog, 3)
	if err != nil {
		return nil, err
	}
	if len(toUsers) == 0 {
		return nil, common.ErrResourceNotFound
	}

	return toUsers, nil
}
func (m *Match) PostNewMatch(fromUserId, toUserId string, matchStatus domain.MatchStatus) (string, error) {
	id, err := m.matchRepo.InsertNewMatch(fromUserId, toUserId, matchStatus)
	if err != nil {
		return "", err
	}
	return id, nil

}
func (m *Match) GetMatchReqToUserId(userId string) ([]domain.MatchUser, error) {
	matcheds, err := m.matchRepo.SelectMatchReqToUserId(userId)
	if err != nil {
		return nil, err
	}

	return matcheds, nil
}

func (m *Match) RequestChange(matchId string, matchStatus domain.MatchStatus) error {
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
	err = m.updateMatch(matchEntity)
	if err != nil {
		return err
	}
	return nil
}

func (m *Match) RevealChange(matchId string, matchStatus domain.MatchStatus) error {
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

	err = m.updateMatch(matchEntity)
	if err != nil {
		return err
	}
	return nil
}

func (m *Match) updateMatch(matchEntity entity.Match) error {
	err := m.matchRepo.UpdateMatchById(matchEntity)
	if err != nil {
		return err
	}
	return nil
}

func (m *Match) getMatchById(matchId string) (entity.Match, error) {
	matchEntity, err := m.matchRepo.GetMatchById(matchId)
	if err != nil {
		return entity.Match{}, err
	}
	return matchEntity, nil
}
