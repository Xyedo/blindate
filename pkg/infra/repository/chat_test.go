package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiError "github.com/xyedo/blindate/pkg/common/error"
	"github.com/xyedo/blindate/pkg/common/util"
	"github.com/xyedo/blindate/pkg/domain/chat"
	chatEntity "github.com/xyedo/blindate/pkg/domain/chat/entities"
	matchEntity "github.com/xyedo/blindate/pkg/domain/match/entities"
	"github.com/xyedo/blindate/pkg/infra/repository"
)

func Test_InsertNewChat(t *testing.T) {
	chatRepo := repository.NewChat(testQuery)
	setup := func(t *testing.T) (convoId, fromUserId, toUserId string) {
		convRepo := repository.NewConversation(testQuery)
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
		require.NoError(t, err)
		convoId, err = convRepo.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)
		return convoId, fromUsr.ID, toUsr.ID
	}
	t.Run("valid new chat", func(t *testing.T) {
		createNewChat(chatRepo, t)
	})
	t.Run("valid new chat w attachment", func(t *testing.T) {
		convoId, fromUserId, _ := setup(t)

		err := chatRepo.InsertNewChat(&chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUserId,
			Messages:       "",
			SentAt:         time.Now(),
			Attachment: &chatEntity.Attachment{
				BlobLink:  util.RandomUUID() + ".ogg",
				MediaType: "application/ogg",
			},
		})
		require.NoError(t, err)
	})
	t.Run("valid new chat w attachment with mp3", func(t *testing.T) {
		convoId, fromUsr, _ := setup(t)

		err := chatRepo.InsertNewChat(&chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUsr,
			Messages:       "",
			SentAt:         time.Now(),
			Attachment: &chatEntity.Attachment{
				BlobLink:  util.RandomUUID() + ".mp3",
				MediaType: "audio/mpeg",
			},
		})
		require.NoError(t, err)
	})
	t.Run("valid chat but invalid attachment type", func(t *testing.T) {
		convoId, fromUsr, _ := setup(t)

		err := chatRepo.InsertNewChat(&chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUsr,
			Messages:       "",
			SentAt:         time.Now(),
			Attachment: &chatEntity.Attachment{
				BlobLink:  util.RandomUUID() + ".png",
				MediaType: "png",
			},
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, apiError.ErrRefNotFound23503)
	})
	t.Run("invalid conversationId", func(t *testing.T) {
		_, fromUsr, _ := setup(t)
		err := chatRepo.InsertNewChat(&chatEntity.DAO{
			ConversationId: util.RandomUUID(),
			Author:         fromUsr,
			Messages:       util.RandomString(12),
			SentAt:         time.Now(),
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, apiError.ErrRefNotFound23503)
	})
	t.Run("invalid author", func(t *testing.T) {
		convoId, _, _ := setup(t)
		err := chatRepo.InsertNewChat(&chatEntity.DAO{
			ConversationId: convoId,
			Author:         util.RandomUUID(),
			Messages:       util.RandomString(12),
			SentAt:         time.Now(),
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, apiError.ErrRefNotFound23503)
	})
	t.Run("invalid reply_to", func(t *testing.T) {
		convoId, fromUsr, _ := setup(t)
		err := chatRepo.InsertNewChat(&chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUsr,
			Messages:       util.RandomString(12),
			ReplyTo: sql.NullString{
				Valid:  true,
				String: util.RandomUUID(),
			},
			SentAt: time.Now(),
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, apiError.ErrRefNotFound23503)
	})
}

func Test_UpdateSeenChatById(t *testing.T) {
	chat := repository.NewChat(testQuery)
	t.Run("valid", func(t *testing.T) {
		conv := repository.NewConversation(testQuery)
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
		require.NoError(t, err)
		convoId, err := conv.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)
		//this should not happen in prod. the author must have link by converstation
		newChat := &chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUsr.ID,
			Messages:       "whatsup sexy!",
			SentAt:         time.Now(),
		}
		err = chat.InsertNewChat(newChat)
		require.NoError(t, err)
		err = chat.InsertNewChat(&chatEntity.DAO{
			ConversationId: convoId,
			Author:         toUsr.ID,
			Messages:       "omg whatsup",
			SentAt:         time.Now(),
		})
		require.NoError(t, err)
		changedChatIds, err := chat.UpdateSeenChat(convoId, fromUsr.ID)
		require.NoError(t, err)
		require.Len(t, changedChatIds, 1)
		assert.NotEmpty(t, changedChatIds[0])
	})
	t.Run("invalid chatId", func(t *testing.T) {
		changedChatIds, err := chat.UpdateSeenChat(util.RandomUUID(), util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, apiError.ErrResourceNotFound)
		assert.Empty(t, changedChatIds)
	})
}

func Test_DeleteChat(t *testing.T) {
	chat := repository.NewChat(testQuery)
	t.Run("valid delete", func(t *testing.T) {
		chatId, _ := createNewChat(chat, t)
		err := chat.DeleteChatById(chatId)
		require.NoError(t, err)
	})
	t.Run("invalid chatId", func(t *testing.T) {
		err := chat.DeleteChatById(util.RandomUUID())
		require.Error(t, err)
		require.ErrorIs(t, err, apiError.ErrRefNotFound23503)
	})
	t.Run("delete refrenced chat", func(t *testing.T) {
		conv := repository.NewConversation(testQuery)
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
		require.NoError(t, err)
		convoId, err := conv.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)
		newChat := &chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUsr.ID,
			Messages:       "whatsup sexy!",
			SentAt:         time.Now(),
		}
		err = chat.InsertNewChat(newChat)
		require.NoError(t, err)
		err = chat.InsertNewChat(&chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUsr.ID,
			Messages:       util.RandomString(12),
			ReplyTo: sql.NullString{
				Valid:  true,
				String: newChat.Id,
			},
			SentAt: time.Now(),
		})
		require.NoError(t, err)
		err = chat.DeleteChatById(newChat.Id)
		require.NoError(t, err)
	})
}

func Test_SelectChat(t *testing.T) {
	chatRepo := repository.NewChat(testQuery)
	conv := repository.NewConversation(testQuery)
	t.Run("valid chat", func(t *testing.T) {
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
		require.NoError(t, err)
		convoId, err := conv.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)
		chatDivider := &chatEntity.DAO{
			ConversationId: convoId,
			Author:         fromUsr.ID,
			Messages:       util.RandomString(15),
			SentAt:         time.Now(),
		}
		for i := 0; i < 20; i++ {
			if i == 10 {
				err = chatRepo.InsertNewChat(chatDivider)
				require.NoError(t, err)
				require.NotEmpty(t, chatDivider.Id)
				err = chatRepo.InsertNewChat(&chatEntity.DAO{
					ConversationId: convoId,
					Author:         toUsr.ID,
					Messages:       util.RandomString(15),
					SentAt:         time.Now(),
				})
				require.NoError(t, err)
				continue
			}
			err = chatRepo.InsertNewChat(&chatEntity.DAO{
				ConversationId: convoId,
				Author:         fromUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
			err = chatRepo.InsertNewChat(&chatEntity.DAO{
				ConversationId: convoId,
				Author:         toUsr.ID,
				Messages:       util.RandomString(15),
				SentAt:         time.Now(),
			})
			require.NoError(t, err)
		}
		t.Run("with default offset", func(t *testing.T) {
			retChats, err := chatRepo.SelectChat(convoId, chat.Filter{})
			require.NoError(t, err)
			require.NotEmpty(t, retChats)
		})
		t.Run("with cursor", func(t *testing.T) {
			retChats, err := chatRepo.SelectChat(convoId, chat.Filter{})
			require.NoError(t, err)
			require.Len(t, retChats, 20)
			before20Chats, err := chatRepo.SelectChat(convoId, chat.Filter{
				Cursor: &chat.Cursor{
					After: true,
					At:    chatDivider.SentAt,
					Id:    chatDivider.Id,
				},
			})
			require.NoError(t, err)
			require.NotEmpty(t, before20Chats)
			after20Chats, err := chatRepo.SelectChat(convoId, chat.Filter{
				Cursor: &chat.Cursor{
					After: false,
					At:    chatDivider.SentAt,
					Id:    chatDivider.Id,
				},
			})
			require.NoError(t, err)
			require.NotEmpty(t, after20Chats)
		})
	})

}
func createNewChat(chat *repository.ChatConn, t *testing.T) (string, string) {
	conv := repository.NewConversation(testQuery)
	matchRepo := repository.NewMatch(testQuery)
	fromUsr := createNewAccount(t)
	toUsr := createNewAccount(t)
	matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchEntity.Requested)
	require.NoError(t, err)
	convoId, err := conv.InsertConversation(matchId)
	require.NoError(t, err)
	require.NotEmpty(t, convoId)
	//this should not happen in prod. the author must have link by converstation
	newChat := &chatEntity.DAO{
		ConversationId: convoId,
		Author:         fromUsr.ID,
		Messages:       "whatsup sexy!",
		SentAt:         time.Now(),
	}
	err = chat.InsertNewChat(newChat)
	require.NoError(t, err)
	return newChat.Id, convoId
}
