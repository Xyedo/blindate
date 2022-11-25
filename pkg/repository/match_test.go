package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/entity"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertNewMatch(t *testing.T) {
	matchRepo := NewMatch(testQuery)
	t.Run("valid newMatch", func(t *testing.T) {
		matchId := createNewMatch(t)
		assert.NotEmpty(t, matchId)
	})
	t.Run("invalid newMatch requestFrom", func(t *testing.T) {
		toUsr := createNewAccount(t)
		_, err := matchRepo.InsertNewMatch(util.RandomUUID(), toUsr.ID)
		require.Error(t, err)
		var pqErr *pq.Error
		require.ErrorAs(t, err, &pqErr)
		assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
		assert.Contains(t, pqErr.Constraint, "match_request_from_fkey")
	})
	t.Run("invalid newMatch requestTo", func(t *testing.T) {
		fromUsr := createNewAccount(t)
		_, err := matchRepo.InsertNewMatch(fromUsr.ID, util.RandomUUID())
		require.Error(t, err)
		var pqErr *pq.Error
		require.ErrorAs(t, err, &pqErr)
		assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
		assert.Contains(t, pqErr.Constraint, "match_request_to_fkey")
	})
	t.Run("double on requestFrom and requestTo", func(t *testing.T) {
		matchRepo := NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)

		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID)
		require.NoError(t, err)
		assert.NotEmpty(t, matchId)
		matchId, err = matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID)
		require.Error(t, err)
		assert.Empty(t, matchId)
		var pqErr *pq.Error
		require.ErrorAs(t, err, &pqErr)
		assert.Equal(t, pq.ErrorCode("23505"), pqErr.Code)
		assert.Contains(t, pqErr.Constraint, "match_request_from_request_to_unique")
	})
}

func Test_SelectMatchByUserId(t *testing.T) {
	matchRepo := NewMatch(testQuery)
	t.Run("valid match", func(t *testing.T) {
		fromUsr := createNewAccount(t)
		for i := 0; i < 5; i++ {
			toUsr := createNewAccount(t)
			matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID)
			require.NoError(t, err)
			assert.NotEmpty(t, matchId)
		}
		matchs, err := matchRepo.SelectMatchByUserId(fromUsr.ID)
		require.NoError(t, err)
		assert.Len(t, matchs, 5)
		assert.Equal(t, fromUsr.ID, matchs[0].RequestFrom)
	})
	t.Run("zero matchs with valid user", func(t *testing.T) {
		user := createNewAccount(t)
		convs, err := matchRepo.SelectMatchByUserId(user.ID)
		require.NoError(t, err)
		assert.Empty(t, convs)
	})
	t.Run("zero matchs with invalid user", func(t *testing.T) {
		convs, err := matchRepo.SelectMatchByUserId(util.RandomUUID())
		require.NoError(t, err)
		assert.Empty(t, convs)
	})
}
func Test_GetMatchById(t *testing.T) {
	matchRepo := NewMatch(testQuery)
	t.Run("valid select", func(t *testing.T) {
		matchId := createNewMatch(t)
		matchRes, err := matchRepo.GetMatchById(matchId)
		require.NoError(t, err)
		assert.Equal(t, matchId, matchRes.Id)
	})
	t.Run("invalid userId", func(t *testing.T) {
		matchRes, err := matchRepo.GetMatchById(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Empty(t, matchRes)
	})
}

func Test_UpdateMatchById(t *testing.T) {
	matchRepo := NewMatch(testQuery)
	setupFunc := func(t *testing.T) entity.Match {
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID)
		require.NoError(t, err)
		require.NotEmpty(t, matchId)
		match, err := matchRepo.GetMatchById(matchId)
		require.NoError(t, err)
		return *match
	}
	t.Run("valid update request_status", func(t *testing.T) {
		newMatch := setupFunc(t)
		newMatch.RequestStatus = string(domain.Requested)
		err := matchRepo.UpdateMatchById(newMatch)
		require.NoError(t, err)
	})
	t.Run("valid update accepted_at", func(t *testing.T) {
		newMatch := setupFunc(t)
		newMatch.AcceptedAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
		err := matchRepo.UpdateMatchById(newMatch)
		require.NoError(t, err)
	})
	t.Run("valid update reveal_status", func(t *testing.T) {
		newMatch := setupFunc(t)
		newMatch.RevealStatus = string(domain.Requested)
		err := matchRepo.UpdateMatchById(newMatch)
		require.NoError(t, err)
	})
	t.Run("valid update revealed_at", func(t *testing.T) {
		newMatch := setupFunc(t)
		newMatch.RevealedAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
		err := matchRepo.UpdateMatchById(newMatch)
		require.NoError(t, err)
	})
	t.Run("invalid matchId", func(t *testing.T) {
		newMatch := setupFunc(t)
		newMatch.Id = util.RandomUUID()
		err := matchRepo.UpdateMatchById(newMatch)
		require.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})

}
func createNewMatch(t *testing.T) string {
	matchRepo := NewMatch(testQuery)
	fromUsr := createNewAccount(t)
	toUsr := createNewAccount(t)

	matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID)
	require.NoError(t, err)
	return matchId
}
