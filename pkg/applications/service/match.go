package service

import (
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain/event"
	"github.com/xyedo/blindate/pkg/domain/location"
	"github.com/xyedo/blindate/pkg/domain/match"
	matchEntity "github.com/xyedo/blindate/pkg/domain/match/entities"
)

func NewMatch(matchRepo match.Repository, locationRepo location.Repository) *Match {
	return &Match{
		matchRepo:    matchRepo,
		locationRepo: locationRepo,
	}
}

type Match struct {
	matchRepo    match.Repository
	locationRepo location.Repository
}

func (m *Match) FindUserToMatch(userId string) ([]matchEntity.UserDTO, error) {
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
func (m *Match) PostNewMatch(fromUserId, toUserId string, matchStatus matchEntity.Status) (string, error) {
	id, err := m.matchRepo.InsertNewMatch(fromUserId, toUserId, matchStatus)
	if err != nil {
		return "", err
	}
	return id, nil

}
func (m *Match) GetMatchReqToUserId(userId string) ([]matchEntity.FullUserDTO, error) {
	matcheds, err := m.matchRepo.SelectMatchReqToUserId(userId)
	if err != nil {
		return nil, err
	}

	return matcheds, nil
}

func (m *Match) RequestChange(matchId string, matchStatus matchEntity.Status) error {
	matchDAO, err := m.GetMatchById(matchId)
	if err != nil {
		return err
	}
	switch matchStatus {
	case matchEntity.Requested:
		if matchDAO.RequestStatus != string(matchEntity.Unknown) {
			return ErrInvalidMatchStatus
		}
		matchDAO.RequestStatus = string(matchEntity.Requested)
	case matchEntity.Declined:
		matchDAO.RequestStatus = string(matchEntity.Declined)
	case matchEntity.Accepted:
		if matchDAO.RequestStatus != string(matchEntity.Requested) {
			return ErrInvalidMatchStatus
		}
		matchDAO.RequestStatus = string(matchEntity.Accepted)
	}
	err = m.updateMatch(matchDAO)
	if err != nil {
		return err
	}
	return nil
}

func (m *Match) RevealChange(matchId string, matchStatus matchEntity.Status) error {
	matchDAO, err := m.GetMatchById(matchId)
	if err != nil {
		return err
	}
	if matchDAO.RequestStatus != string(matchEntity.Accepted) {
		return ErrInvalidMatchStatus
	}
	switch matchStatus {
	case matchEntity.Requested:
		if matchDAO.RevealStatus != string(matchEntity.Unknown) {
			return ErrInvalidMatchStatus
		}
		matchDAO.RevealStatus = string(matchEntity.Requested)
		event.MatchRevealed.Trigger(event.MatchRevealedPayload{
			MatchId:     matchId,
			MatchStatus: matchEntity.Requested,
		})
	case matchEntity.Declined:
		if matchDAO.RevealStatus == string(matchEntity.Unknown) {
			return ErrInvalidMatchStatus
		}
		matchDAO.RevealStatus = string(matchEntity.Declined)
		event.MatchRevealed.Trigger(event.MatchRevealedPayload{
			MatchId:     matchId,
			MatchStatus: matchEntity.Declined,
		})
	case matchEntity.Accepted:
		if matchDAO.RevealStatus != string(matchEntity.Requested) {
			return ErrInvalidMatchStatus
		}
		matchDAO.RevealStatus = string(matchEntity.Accepted)
		event.MatchRevealed.Trigger(event.MatchRevealedPayload{
			MatchId:     matchId,
			MatchStatus: matchEntity.Accepted,
		})
	}

	err = m.updateMatch(matchDAO)
	if err != nil {
		return err
	}
	return nil
}

func (m *Match) GetMatchById(matchId string) (matchEntity.MatchDAO, error) {
	matchDAO, err := m.matchRepo.GetMatchById(matchId)
	if err != nil {
		return matchEntity.MatchDAO{}, err
	}
	return matchDAO, nil
}

func (m *Match) updateMatch(matchEntity matchEntity.MatchDAO) error {
	err := m.matchRepo.UpdateMatchById(matchEntity)
	if err != nil {
		return err
	}
	return nil
}
