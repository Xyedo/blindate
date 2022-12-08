package service

import (
	"errors"
	"fmt"

	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
)

var (
	ErrRefMatchId         = fmt.Errorf("%w:invalid matchId", common.ErrRefNotFound23503)
	ErrInvalidMatchStatus = errors.New("not yet accepted/revealed in matchId")
)

func NewConversation(convRepo repository.Conversation, matchRepo repository.Match) *Conversation {
	return &Conversation{
		convRepo:  convRepo,
		matchRepo: matchRepo,
	}
}

type Conversation struct {
	convRepo  repository.Conversation
	matchRepo repository.Match
}

func (c *Conversation) CreateConversation(matchId string) (string, error) {
	matchEntity, err := c.matchRepo.GetMatchById(matchId)
	if err != nil {
		return "", err
	}
	if matchEntity.RequestStatus != string(domain.Accepted) {
		return "", ErrInvalidMatchStatus
	}
	id, err := c.convRepo.InsertConversation(matchId)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (c *Conversation) FindConversationById(matchId string) (domain.Conversation, error) {
	conv, err := c.convRepo.SelectConversationById(matchId)
	if err != nil {
		return domain.Conversation{}, err
	}

	if conv.RevealStatus != string(domain.Accepted) {
		conv.FromUser.FullName = ""
		conv.FromUser.ProfilePic = ""
		conv.ToUser.FullName = ""
		conv.ToUser.ProfilePic = ""
	}

	return conv, nil
}

func (c *Conversation) GetConversationByUserId(userId string) ([]domain.Conversation, error) {
	convs, err := c.convRepo.SelectConversationByUserId(userId, nil)
	if err != nil {
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
func (c *Conversation) DeleteConversationById(convoId string) error {
	err := c.convRepo.DeleteConversationById(convoId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conversation) UpdateConvRow(convoId string) error {
	err := c.convRepo.UpdateChatRow(convoId)
	if err != nil {
		return err
	}
	return nil
}

func (c *Conversation) UpdateConvDay(convoId string) error {
	err := c.convRepo.UpdateDayPass(convoId)
	if err != nil {
		return err
	}
	return nil
}
