package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

var (
	ErrConvoWithSelf  = fmt.Errorf("conversation: cannot create conversation with yourself")
	ErrRefMatchId     = fmt.Errorf("%w:invalid matchId", domain.ErrRefNotFound23503)
	ErrNotYetAccepted = errors.New("not yet accepted/revealed in matchId")
)

func NewConversation(convRepo repository.Conversation, matchRepo repository.Match) *conversation {
	return &conversation{
		convRepo:  convRepo,
		matchRepo: matchRepo,
	}
}

type conversation struct {
	convRepo  repository.Conversation
	matchRepo repository.Match
}

func (c *conversation) CreateConversation(matchId string) (string, error) {
	matchEntity, err := c.matchRepo.GetMatchById(matchId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", domain.ErrTooLongAccesingDB
		}
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrRefMatchId
		}
		return "", err
	}
	if matchEntity.RequestStatus != string(domain.Accepted) {
		return "", ErrNotYetAccepted
	}
	id, err := c.convRepo.InsertConversation(matchId)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return "", domain.ErrTooLongAccesingDB
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23503" {
				return "", domain.ErrRefNotFound23503
			}
			if pqErr.Code == "23505" {
				return "", domain.ErrUniqueConstraint23505
			}
			return "", pqErr
		}
		return "", err
	}
	return id, nil
}

func (c *conversation) FindConversationById(matchId string) (*domain.Conversation, error) {
	conv, err := c.convRepo.SelectConversationById(matchId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccesingDB
		}
		return nil, err
	}

	if conv.RevealStatus != string(domain.Accepted) {
		conv.FromUser.FullName = ""
		conv.FromUser.ProfilePic = ""
		conv.ToUser.FullName = ""
		conv.ToUser.ProfilePic = ""
	}

	return conv, nil
}

func (c *conversation) GetConversationByUserId(userId string) ([]domain.Conversation, error) {
	convs, err := c.convRepo.SelectConversationByUserId(userId, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return nil, domain.ErrTooLongAccesingDB
		}
		return nil, err
	}
	for i := range convs {
		if convs[i].RevealStatus != string(domain.Accepted) {
			convs[i].FromUser.FullName = ""
			convs[i].FromUser.ProfilePic = ""
			convs[i].ToUser.FullName = ""
			convs[i].ToUser.ProfilePic = ""
		}
	}
	return convs, nil
}
func (c *conversation) DeleteConversationById(convoId string) error {
	err := c.convRepo.DeleteConversationById(convoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	return nil
}

func (c *conversation) UpdateConvRow(convoId string) error {
	err := c.convRepo.UpdateChatRow(convoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	return nil
}

func (c *conversation) UpdateConvDay(convoId string) error {
	err := c.convRepo.UpdateDayPass(convoId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrResourceNotFound
		}
		if errors.Is(err, context.Canceled) {
			return domain.ErrTooLongAccesingDB
		}
		return err
	}
	return nil
}
