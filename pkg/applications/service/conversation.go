package service

import (
	"errors"
	"fmt"

	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain/conversation"
	convEntity "github.com/xyedo/blindate/pkg/domain/conversation/entities"
	"github.com/xyedo/blindate/pkg/domain/match"
	matchEntity "github.com/xyedo/blindate/pkg/domain/match/entities"
)

var (
	ErrRefMatchId         = fmt.Errorf("%w:invalid matchId", common.ErrRefNotFound23503)
	ErrInvalidMatchStatus = errors.New("not yet accepted/revealed in matchId")
)

func NewConversation(convRepo conversation.Repository, matchRepo match.Repository) *Conversation {
	return &Conversation{
		convRepo:  convRepo,
		matchRepo: matchRepo,
	}
}

type Conversation struct {
	convRepo  conversation.Repository
	matchRepo match.Repository
}

func (c *Conversation) CreateConversation(matchId string) (string, error) {
	matchDAO, err := c.matchRepo.GetMatchById(matchId)
	if err != nil {
		return "", err
	}
	if matchDAO.RequestStatus != string(matchEntity.Accepted) {
		return "", ErrInvalidMatchStatus
	}
	id, err := c.convRepo.InsertConversation(matchId)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (c *Conversation) FindConversationById(matchId string) (convEntity.DTO, error) {
	conv, err := c.convRepo.SelectConversationById(matchId)
	if err != nil {
		return convEntity.DTO{}, err
	}

	if conv.RevealStatus != string(matchEntity.Accepted) {
		conv.FromUser.FullName = ""
		conv.FromUser.ProfilePic = ""
		conv.ToUser.FullName = ""
		conv.ToUser.ProfilePic = ""
	}

	return conv, nil
}

func (c *Conversation) GetConversationByUserId(userId string) ([]convEntity.DTO, error) {
	convs, err := c.convRepo.SelectConversationByUserId(userId, nil)
	if err != nil {
		return nil, err
	}
	for i := range convs {
		if convs[i].RevealStatus != string(matchEntity.Accepted) {
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
