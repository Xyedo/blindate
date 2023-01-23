package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain/chat"
	chatEntity "github.com/xyedo/blindate/pkg/domain/chat/entities"
	"github.com/xyedo/blindate/pkg/domain/conversation"
	matchEntity "github.com/xyedo/blindate/pkg/domain/match/entities"
	"github.com/xyedo/blindate/pkg/infra/repository"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertConversation(t *testing.T) {
	conv := repository.NewConversation(testQuery)
	t.Run("valid new conv", func(t *testing.T) {
		createNewConvo(conv, t)
	})
	t.Run("invalid new conv", func(t *testing.T) {
		id, err := conv.InsertConversation(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
		require.Empty(t, id)

	})
	t.Run("duplicate conv", func(t *testing.T) {
		matchId := createNewMatch(t)
		id, err := conv.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, id)
		id, err = conv.InsertConversation(matchId)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)
		require.Empty(t, id)

	})
}

func Test_SelectConversationById(t *testing.T) {
	conv := repository.NewConversation(testQuery)
	matchRepo := repository.NewMatch(testQuery)
	chat := repository.NewChat(testQuery)
	user := repository.NewUser(testQuery)
	t.Run("valid select with full attributes", func(t *testing.T) {
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		var expectedProfilePicCreator string
		for i := 0; i < 4; i++ {
			if i == 2 {
				expectedProfilePicCreator = "true.png"
				_, err := user.CreateProfilePicture(fromUsr.ID, expectedProfilePicCreator, true)
				require.NoError(t, err)
				continue
			}
			_, err := user.CreateProfilePicture(fromUsr.ID, fmt.Sprintf("%d.png", i), false)
			require.NoError(t, err)
		}
		var expectedProfilePicRecipient string
		for i := 0; i < 4; i++ {
			if i == 3 {
				expectedProfilePicRecipient = fmt.Sprintf("%d.png", i)
				_, err := user.CreateProfilePicture(toUsr.ID, expectedProfilePicRecipient, false)
				require.NoError(t, err)
				continue
			}
			_, err := user.CreateProfilePicture(toUsr.ID, fmt.Sprintf("%d.png", i), false)
			require.NoError(t, err)
		}
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
		require.NoError(t, err)
		convoId, err := conv.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)

		for i := 0; i < 4; i++ {
			err = chat.InsertNewChat(&chatEntity.DAO{
				ConversationId: convoId,
				Author:         fromUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
			err = chat.InsertNewChat(&chatEntity.DAO{
				ConversationId: convoId,
				Author:         toUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
		}
		expectedCreatorMsg := "hey sexy!"
		err = chat.InsertNewChat(&chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUsr.ID,
			Messages:       expectedCreatorMsg,
			SentAt:         time.Now(),
		})
		require.NoError(t, err)

		conv, err := conv.SelectConversationById(convoId)
		require.NoError(t, err)
		require.NotEmpty(t, conv)

		require.NotEmpty(t, expectedProfilePicCreator)
		require.NotEmpty(t, expectedProfilePicRecipient)

		assert.Equal(t, expectedProfilePicCreator, conv.FromUser.ProfilePic)
		assert.Equal(t, expectedProfilePicRecipient, conv.ToUser.ProfilePic)
		assert.Equal(t, expectedCreatorMsg, conv.LastMessage)
		assert.Equal(t, "requested", conv.RequestStatus)
		assert.Equal(t, "unknown", conv.RevealStatus)
	})
	t.Run("invalid convoId", func(t *testing.T) {
		conv, err := conv.SelectConversationById(util.RandomUUID())
		require.Error(t, err)
		require.Empty(t, conv)
		assert.ErrorIs(t, err, common.ErrResourceNotFound)

	})
}

func Test_SelectConvoByUserId(t *testing.T) {
	conv := repository.NewConversation(testQuery)
	matchRepo := repository.NewMatch(testQuery)
	chat := repository.NewChat(testQuery)
	user := repository.NewUser(testQuery)
	t.Run("valid select with full attributes", func(t *testing.T) {
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		var expectedProfilePicCreator string
		for i := 0; i < 4; i++ {
			if i == 2 {
				expectedProfilePicCreator = "true.png"
				_, err := user.CreateProfilePicture(fromUsr.ID, expectedProfilePicCreator, true)
				require.NoError(t, err)
				continue
			}
			_, err := user.CreateProfilePicture(fromUsr.ID, fmt.Sprintf("%d.png", i), false)
			require.NoError(t, err)
		}
		var expectedProfilePicRecipient string
		for i := 0; i < 4; i++ {
			if i == 3 {
				expectedProfilePicRecipient = fmt.Sprintf("%d.png", i)
				_, err := user.CreateProfilePicture(toUsr.ID, expectedProfilePicRecipient, false)
				require.NoError(t, err)
				continue
			}
			_, err := user.CreateProfilePicture(toUsr.ID, fmt.Sprintf("%d.png", i), false)
			require.NoError(t, err)
		}
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
		require.NoError(t, err)
		convoId, err := conv.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)

		for i := 0; i < 4; i++ {
			err = chat.InsertNewChat(&chatEntity.DAO{
				ConversationId: convoId,
				Author:         fromUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
			err = chat.InsertNewChat(&chatEntity.DAO{
				ConversationId: convoId,
				Author:         toUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
		}
		expectedCreatorMsg := "hey sexy!"
		err = chat.InsertNewChat(&chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUsr.ID,
			Messages:       expectedCreatorMsg,
			SentAt:         time.Now(),
		})
		require.NoError(t, err)

		res, err := conv.SelectConversationByUserId(fromUsr.ID, nil)
		require.NoError(t, err)
		require.Len(t, res, 1)

		require.NotEmpty(t, expectedProfilePicCreator)
		require.NotEmpty(t, expectedProfilePicRecipient)

		assert.Equal(t, expectedProfilePicCreator, res[0].FromUser.ProfilePic)
		assert.Equal(t, expectedProfilePicRecipient, res[0].ToUser.ProfilePic)
		assert.Equal(t, expectedCreatorMsg, res[0].LastMessage)

	})
	t.Run("valid select with little to none attributes", func(t *testing.T) {
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
		require.NoError(t, err)
		_, err = conv.InsertConversation(matchId)
		require.NoError(t, err)

		res, err := conv.SelectConversationByUserId(fromUsr.ID, nil)
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, "", res[0].FromUser.ProfilePic)
		assert.Equal(t, "", res[0].ToUser.ProfilePic)
		assert.Equal(t, "", res[0].LastMessage)
	})
	t.Run("valid with full attr and lot match", func(t *testing.T) {
		fromUsr := createNewAccount(t)
		_, err := user.CreateProfilePicture(fromUsr.ID, util.RandomUUID()+".png", true)
		require.NoError(t, err)
		for i := 0; i < 30; i++ {
			toUsr := createNewAccount(t)
			_, err = user.CreateProfilePicture(toUsr.ID, util.RandomUUID()+".png", false)
			require.NoError(t, err)
			matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
			require.NoError(t, err)
			convoId, err := conv.InsertConversation(matchId)
			require.NoError(t, err)

			if util.RandomBool() {
				newChat := &chatEntity.DAO{
					ConversationId: convoId,
					Messages:       util.RandomString(10),
					SentAt:         time.Now(),
				}
				if util.RandomBool() {
					newChat.Author = fromUsr.ID
				} else {
					newChat.Author = toUsr.ID
				}
				err = chat.InsertNewChat(newChat)
				require.NoError(t, err)
			}
		}
		t.Run("without offset", func(t *testing.T) {
			res, err := conv.SelectConversationByUserId(fromUsr.ID, nil)
			require.NoError(t, err)
			require.NotEmpty(t, res)
			require.Len(t, res, 20)
		})
		t.Run("with offset", func(t *testing.T) {
			res, err := conv.SelectConversationByUserId(fromUsr.ID, &conversation.Filter{Offset: 10})
			require.NoError(t, err)
			require.NotEmpty(t, res)
			require.Len(t, res, 20)
		})

	})
	t.Run("return 0 length 'coz userId didnt exists", func(t *testing.T) {
		res, err := conv.SelectConversationByUserId(util.RandomUUID(), nil)
		require.NoError(t, err)
		assert.Empty(t, res)
	})
}

func Test_UpdateChatRow(t *testing.T) {
	conv := repository.NewConversation(testQuery)
	t.Run("valid update", func(t *testing.T) {
		convoId := createNewConvo(conv, t)
		err := conv.UpdateChatRow(convoId)
		require.NoError(t, err)
	})
	t.Run("invalid convoId", func(t *testing.T) {
		err := conv.UpdateChatRow(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
}
func Test_UpdateDayPass(t *testing.T) {
	conv := repository.NewConversation(testQuery)
	t.Run("valid update", func(t *testing.T) {
		convoId := createNewConvo(conv, t)
		err := conv.UpdateDayPass(convoId)
		require.NoError(t, err)
	})
	t.Run("invalid convoId", func(t *testing.T) {
		err := conv.UpdateDayPass(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
}

func Test_DeleteConvoById(t *testing.T) {
	conv := repository.NewConversation(testQuery)

	t.Run("valid", func(t *testing.T) {
		chatRepo := repository.NewChat(testQuery)
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
		require.NoError(t, err)
		convoId, err := conv.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)
		for i := 0; i < 5; i++ {
			newChat := &chatEntity.DAO{
				ConversationId: convoId,
				Messages:       util.RandomString(12),
				SentAt:         time.Now(),
			}
			if util.RandomBool() {
				newChat.Author = fromUsr.ID
			} else {
				newChat.Author = toUsr.ID
			}
			err := chatRepo.InsertNewChat(newChat)
			require.NoError(t, err)
		}
		err = conv.DeleteConversationById(convoId)
		require.NoError(t, err)
		convs, err := conv.SelectConversationById(convoId)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrResourceNotFound)
		assert.Empty(t, convs)
		chats, err := chatRepo.SelectChat(convoId, chat.Filter{
			Limit: 10,
		})
		require.NoError(t, err)
		assert.Empty(t, chats)

	})
	t.Run("invalid convoId", func(t *testing.T) {
		err := conv.DeleteConversationById(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
}
func createNewConvo(conv *repository.ConvConn, t *testing.T) string {
	matchId := createNewMatch(t)
	id, err := conv.InsertConversation(matchId)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	return id
}
