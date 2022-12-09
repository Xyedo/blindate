package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/domain/entity"
	"github.com/xyedo/blindate/pkg/repository"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertNewChat(t *testing.T) {
	chat := repository.NewChat(testQuery)
	setup := func(t *testing.T) (convoId, fromUserId, toUserId string) {
		convRepo := repository.NewConversation(testQuery)
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, domain.Requested)
		require.NoError(t, err)
		convoId, err = convRepo.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)
		return convoId, fromUsr.ID, toUsr.ID
	}
	t.Run("valid new chat", func(t *testing.T) {
		createNewChat(chat, t)
	})
	t.Run("valid new chat w attachment", func(t *testing.T) {
		convoId, fromUserId, _ := setup(t)

		err := chat.InsertNewChat(&entity.Chat{
			ConversationId: convoId,
			Author:         fromUserId,
			Messages:       "",
			SentAt:         time.Now(),
			Attachment: &domain.ChatAttachment{
				BlobLink:  util.RandomUUID() + ".ogg",
				MediaType: "application/ogg",
			},
		})
		require.NoError(t, err)
	})
	t.Run("valid new chat w attachment with mp3", func(t *testing.T) {
		convoId, fromUsr, _ := setup(t)

		err := chat.InsertNewChat(&entity.Chat{
			ConversationId: convoId,
			Author:         fromUsr,
			Messages:       "",
			SentAt:         time.Now(),
			Attachment: &domain.ChatAttachment{
				BlobLink:  util.RandomUUID() + ".mp3",
				MediaType: "audio/mpeg",
			},
		})
		require.NoError(t, err)
	})
	t.Run("valid chat but invalid attachment type", func(t *testing.T) {
		convoId, fromUsr, _ := setup(t)

		err := chat.InsertNewChat(&entity.Chat{
			ConversationId: convoId,
			Author:         fromUsr,
			Messages:       "",
			SentAt:         time.Now(),
			Attachment: &domain.ChatAttachment{
				BlobLink:  util.RandomUUID() + ".png",
				MediaType: "png",
			},
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("invalid conversationId", func(t *testing.T) {
		_, fromUsr, _ := setup(t)
		err := chat.InsertNewChat(&entity.Chat{
			ConversationId: util.RandomUUID(),
			Author:         fromUsr,
			Messages:       util.RandomString(12),
			SentAt:         time.Now(),
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("invalid author", func(t *testing.T) {
		convoId, _, _ := setup(t)
		err := chat.InsertNewChat(&entity.Chat{
			ConversationId: convoId,
			Author:         util.RandomUUID(),
			Messages:       util.RandomString(12),
			SentAt:         time.Now(),
		})
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("invalid reply_to", func(t *testing.T) {
		convoId, fromUsr, _ := setup(t)
		err := chat.InsertNewChat(&entity.Chat{
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
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
}

func Test_UpdateSeenChatById(t *testing.T) {
	chat := repository.NewChat(testQuery)
	t.Run("valid", func(t *testing.T) {
		chatId := createNewChat(chat, t)
		err := chat.UpdateSeenChatById(chatId)
		require.NoError(t, err)
	})
	t.Run("invalid chatId", func(t *testing.T) {
		err := chat.UpdateSeenChatById(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
}

func Test_DeleteChat(t *testing.T) {
	chat := repository.NewChat(testQuery)
	t.Run("valid delete", func(t *testing.T) {
		chatId := createNewChat(chat, t)
		err := chat.DeleteChatById(chatId)
		require.NoError(t, err)
	})
	t.Run("invalid chatId", func(t *testing.T) {
		err := chat.DeleteChatById(util.RandomUUID())
		require.Error(t, err)
		require.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("delete refrenced chat", func(t *testing.T) {
		conv := repository.NewConversation(testQuery)
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, domain.Requested)
		require.NoError(t, err)
		convoId, err := conv.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)
		newChat := &entity.Chat{
			ConversationId: convoId,
			Author:         fromUsr.ID,
			Messages:       "whatsup sexy!",
			SentAt:         time.Now(),
		}
		err = chat.InsertNewChat(newChat)
		require.NoError(t, err)
		err = chat.InsertNewChat(&entity.Chat{
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
	chat := repository.NewChat(testQuery)
	conv := repository.NewConversation(testQuery)
	t.Run("valid chat", func(t *testing.T) {
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, domain.Requested)
		require.NoError(t, err)
		convoId, err := conv.InsertConversation(matchId)
		require.NoError(t, err)
		require.NotEmpty(t, convoId)
		chatDivider := &entity.Chat{
			ConversationId: convoId,
			Author:         fromUsr.ID,
			Messages:       util.RandomString(15),
			SentAt:         time.Now(),
		}
		for i := 0; i < 20; i++ {
			if i == 10 {
				err = chat.InsertNewChat(chatDivider)
				require.NoError(t, err)
				require.NotEmpty(t, chatDivider.Id)
				err = chat.InsertNewChat(&entity.Chat{
					ConversationId: convoId,
					Author:         toUsr.ID,
					Messages:       util.RandomString(15),
					SentAt:         time.Now(),
				})
				require.NoError(t, err)
				continue
			}
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
		t.Run("with default offset", func(t *testing.T) {
			retChats, err := chat.SelectChat(convoId, entity.ChatFilter{})
			require.NoError(t, err)
			require.NotEmpty(t, retChats)
		})
		t.Run("with cursor", func(t *testing.T) {
			retChats, err := chat.SelectChat(convoId, entity.ChatFilter{})
			require.NoError(t, err)
			require.Len(t, retChats, 20)
			before20Chats, err := chat.SelectChat(convoId, entity.ChatFilter{
				Cursor: &entity.ChatCursor{
					After: true,
					At:    chatDivider.SentAt,
					Id:    chatDivider.Id,
				},
			})
			require.NoError(t, err)
			require.NotEmpty(t, before20Chats)
			after20Chats, err := chat.SelectChat(convoId, entity.ChatFilter{
				Cursor: &entity.ChatCursor{
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
func createNewChat(chat *repository.ChatConn, t *testing.T) string {
	conv := repository.NewConversation(testQuery)
	matchRepo := repository.NewMatch(testQuery)
	fromUsr := createNewAccount(t)
	toUsr := createNewAccount(t)
	matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, domain.Requested)
	require.NoError(t, err)
	convoId, err := conv.InsertConversation(matchId)
	require.NoError(t, err)
	require.NotEmpty(t, convoId)
	//this should not happen in prod. the author must have link by converstation
	newChat := &entity.Chat{
		ConversationId: convoId,
		Author:         fromUsr.ID,
		Messages:       "whatsup sexy!",
		SentAt:         time.Now(),
	}
	err = chat.InsertNewChat(newChat)
	require.NoError(t, err)
	return newChat.Id
}
