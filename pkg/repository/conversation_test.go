package repository

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/entity"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertConversation(t *testing.T) {
	conv := NewConversation(testQuery)
	t.Run("valid new conv", func(t *testing.T) {
		createNewConvo(conv, t)
	})
	t.Run("invalid new conv", func(t *testing.T) {
		fromUsr := createNewAccount(t)
		id, err := conv.InsertConversation(fromUsr.ID, util.RandomUUID())
		require.Error(t, err)
		require.Empty(t, id)
		var pqErr *pq.Error
		require.ErrorAs(t, err, &pqErr)
		require.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
		require.Contains(t, pqErr.Constraint, "to_id")
	})
	t.Run("duplicate conv", func(t *testing.T) {
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		id, err := conv.InsertConversation(fromUsr.ID, toUsr.ID)
		require.NoError(t, err)
		require.NotEmpty(t, id)
		id, err = conv.InsertConversation(fromUsr.ID, toUsr.ID)
		require.Error(t, err)
		require.Empty(t, id)
		var pqErr *pq.Error
		require.ErrorAs(t, err, &pqErr)
		assert.Equal(t, pq.ErrorCode("23505"), pqErr.Code)
		assert.Contains(t, pqErr.Constraint, "to_id")
	})
}

func Test_SelectConversationById(t *testing.T) {
	conv := NewConversation(testQuery)
	chat := NewChat(testQuery)
	user := NewUser(testQuery)
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

		convoId, err := conv.InsertConversation(fromUsr.ID, toUsr.ID)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)

		for i := 0; i < 4; i++ {
			err = chat.InsertNewChat(&entity.Chat{
				ConversationId: convoId,
				Author:         fromUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
			err = chat.InsertNewChat(&entity.Chat{
				ConversationId: convoId,
				Author:         toUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
		}
		expectedCreatorMsg := "hey sexy!"
		err = chat.InsertNewChat(&entity.Chat{
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
	})
	t.Run("invalid convoId", func(t *testing.T) {
		conv, err := conv.SelectConversationById(util.RandomUUID())
		require.Error(t, err)
		require.Empty(t, conv)
		assert.ErrorIs(t, err, sql.ErrNoRows)

	})
}

func Test_SelectConvoByUserId(t *testing.T) {
	conv := NewConversation(testQuery)
	chat := NewChat(testQuery)
	user := NewUser(testQuery)
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

		convoId, err := conv.InsertConversation(fromUsr.ID, toUsr.ID)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)

		for i := 0; i < 4; i++ {
			err = chat.InsertNewChat(&entity.Chat{
				ConversationId: convoId,
				Author:         fromUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
			err = chat.InsertNewChat(&entity.Chat{
				ConversationId: convoId,
				Author:         toUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
		}
		expectedCreatorMsg := "hey sexy!"
		err = chat.InsertNewChat(&entity.Chat{
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
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		_, err := conv.InsertConversation(fromUsr.ID, toUsr.ID)
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
			convoId, err := conv.InsertConversation(fromUsr.ID, toUsr.ID)
			require.NoError(t, err)

			if util.RandomBool() {
				newChat := &entity.Chat{
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
			res, err := conv.SelectConversationByUserId(fromUsr.ID, &entity.ConvFilter{Offset: 10})
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
	conv := NewConversation(testQuery)
	t.Run("valid update", func(t *testing.T) {
		convoId := createNewConvo(conv, t)
		err := conv.UpdateChatRow(convoId)
		require.NoError(t, err)
	})
	t.Run("invalid convoId", func(t *testing.T) {
		err := conv.UpdateChatRow(util.RandomUUID())
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
func Test_UpdateDayPass(t *testing.T) {
	conv := NewConversation(testQuery)
	t.Run("valid update", func(t *testing.T) {
		convoId := createNewConvo(conv, t)
		err := conv.UpdateDayPass(convoId)
		require.NoError(t, err)
	})
	t.Run("invalid convoId", func(t *testing.T) {
		err := conv.UpdateDayPass(util.RandomUUID())
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})
}
func Test_DeleteConvoById(t *testing.T) {
	conv := NewConversation(testQuery)
	t.Run("valid", func(t *testing.T) {
		convoId := createNewConvo(conv, t)
		err := conv.DeleteConversationById(convoId)
		require.NoError(t, err)
	})
	t.Run("invalid convoId", func(t *testing.T) {
		err := conv.DeleteConversationById(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
func createNewConvo(conv *conversation, t *testing.T) string {
	fromUsr := createNewAccount(t)
	toUsr := createNewAccount(t)
	id, err := conv.InsertConversation(fromUsr.ID, toUsr.ID)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	return id
}
