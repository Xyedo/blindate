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

func (m *match) FindNewMatch(userId string) ([]domain.User, error) {
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
	for _, toUser := range toUsers {
		_, err := m.matchRepo.InsertNewMatch(userId, toUser.ID)
		if err != nil {
			var pqErr *pq.Error
			switch {
			case errors.Is(err, context.Canceled):
				return nil, domain.ErrTooLongAccessingDB
			case errors.As(err, &pqErr):
				switch {
				case pqErr.Code == "23503":
					return nil, domain.ErrRefNotFound23503
				case pqErr.Code == "23505":
					return nil, domain.ErrUniqueConstraint23505
				default:
					return nil, pqErr
				}
			default:
				return nil, err
			}
		}
	}
	return toUsers, nil
}
func (m *match) GetMatchByUserId(userId string) ([]domain.Match, error) {
	matchsEn, err := m.matchRepo.SelectMatchByUserId(userId)
	if err != nil {
		return nil, err
	}
	matchs := make([]domain.Match, 0, len(matchsEn))
	for _, matchEn := range matchsEn {
		matchs = append(matchs, *m.convertToDomain(&matchEn))
	}
	return matchs, nil
}

func (m *match) AcceptRequest(matchId string) error {
	matchEntity, err := m.getMatchById(matchId)
	if err != nil {
		return err
	}
	matchEntity.RequestStatus = string(domain.Accepted)

	err = m.updateMatch(*matchEntity)
	if err != nil {
		return err
	}
	return nil
}

func (m *match) RequestReveal(matchId string) error {
	matchEntity, err := m.getMatchById(matchId)
	if err != nil {
		return err
	}
	if matchEntity.RequestStatus != string(domain.Accepted) {
		return ErrNotYetAccepted
	}
	matchEntity.RevealStatus = string(domain.Requested)

	err = m.updateMatch(*matchEntity)
	if err != nil {
		return err
	}
	return nil
}

func (m *match) AcceptReveal(matchId string) error {
	matchEntity, err := m.getMatchById(matchId)
	if err != nil {
		return err
	}
	if !(matchEntity.RequestStatus == string(domain.Accepted) && matchEntity.RevealStatus == string(domain.Requested)) {
		return ErrNotYetAccepted
	}
	matchEntity.RevealStatus = string(domain.Accepted)

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

//	func (*match) convertToEntity(matchDto *domain.Match) *entity.Match {
//		matchEntity := &entity.Match{
//			Id:            matchDto.Id,
//			RequestFrom:   matchDto.RequestFrom,
//			RequestTo:     matchDto.RequestTo,
//			RequestStatus: string(matchDto.RequestStatus),
//			CreatedAt:     matchDto.CreatedAt,
//			RevealStatus:  string(matchDto.RevealStatus),
//		}
//		if matchDto.AcceptedAt != nil {
//			matchEntity.AcceptedAt = sql.NullTime{
//				Valid: true,
//				Time:  *matchDto.AcceptedAt,
//			}
//		}
//		if matchDto.RevealedAt != nil {
//			matchEntity.RevealedAt = sql.NullTime{
//				Valid: true,
//				Time:  *matchDto.RevealedAt,
//			}
//		}
//		return matchEntity
//	}
func (*match) convertToDomain(matchEntity *entity.Match) *domain.Match {
	matchDto := &domain.Match{
		Id:            matchEntity.Id,
		RequestFrom:   matchEntity.RequestFrom,
		RequestTo:     matchEntity.RequestTo,
		RequestStatus: domain.MatchStatus(matchEntity.RequestStatus),
		CreatedAt:     matchEntity.CreatedAt,
		RevealStatus:  domain.MatchStatus(matchEntity.RevealStatus),
	}
	if matchEntity.AcceptedAt.Valid {
		matchDto.AcceptedAt = &matchEntity.AcceptedAt.Time
	}
	if matchEntity.RevealedAt.Valid {
		matchDto.RevealedAt = &matchEntity.RevealedAt.Time
	}
	return matchDto
}
