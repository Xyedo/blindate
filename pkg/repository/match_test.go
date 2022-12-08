package repository_test

import (
	"database/sql"
	"encoding/json"
	"log"
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

func Test_InsertNewMatch(t *testing.T) {
	matchRepo := repository.NewMatch(testQuery)
	t.Run("valid newMatch", func(t *testing.T) {
		matchId := createNewMatch(t)
		require.NotEmpty(t, matchId)
	})
	t.Run("invalid newMatch requestFrom", func(t *testing.T) {
		toUsr := createNewAccount(t)
		_, err := matchRepo.InsertNewMatch(util.RandomUUID(), toUsr.ID, domain.Accepted)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("invalid newMatch requestTo", func(t *testing.T) {
		fromUsr := createNewAccount(t)
		_, err := matchRepo.InsertNewMatch(fromUsr.ID, util.RandomUUID(), domain.Declined)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)

	})
	t.Run("double on requestFrom and requestTo", func(t *testing.T) {
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)

		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, domain.Unknown)
		require.NoError(t, err)
		assert.NotEmpty(t, matchId)
		matchId, err = matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, domain.Accepted)
		require.Error(t, err)
		assert.Empty(t, matchId)
		assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)
	})
	t.Run("invalid requestTo", func(t *testing.T) {
		matchRepo := repository.NewMatch(testQuery)
		fromUsr := createNewAccount(t)

		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, util.RandomUUID(), domain.Unknown)
		require.Error(t, err)
		require.Zero(t, matchId)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("invalid requestFrom", func(t *testing.T) {
		matchRepo := repository.NewMatch(testQuery)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(util.RandomUUID(), toUsr.ID, domain.Unknown)
		require.Error(t, err)
		require.Zero(t, matchId)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
}

func Test_SelectMatchReqToUserId(t *testing.T) {
	matchRepo := repository.NewMatch(testQuery)
	t.Run("valid match", func(t *testing.T) {
		toUser := createNewAccount(t)
		var ExpectedfirstFirstUserId string
		for i := 0; i < 5; i++ {
			fromUsr := createNewAccount(t)
			if i == 0 {
				ExpectedfirstFirstUserId = fromUsr.ID
			}
			bio := createNewInterestBio(t, fromUsr.ID)
			createNewInterestHobbie(t, bio.Id)
			createNewInterestMovieSeries(t, bio.Id)
			createNewInterestSport(t, bio.Id)
			createNewInterestTraveling(t, bio.Id)
			matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUser.ID, domain.Requested)
			require.NoError(t, err)
			assert.NotEmpty(t, matchId)
			if i%2 == 0 {
				fromUsrOdd := createNewAccount(t)
				matchId, err := matchRepo.InsertNewMatch(fromUsrOdd.ID, toUser.ID, domain.Unknown)
				require.NoError(t, err)
				assert.NotEmpty(t, matchId)
			}
		}

		matchs, err := matchRepo.SelectMatchReqToUserId(toUser.ID)
		require.NoError(t, err)
		require.NotEmpty(t, matchs)
		assert.Equal(t, ExpectedfirstFirstUserId, matchs[0].UserId)
		jsonCandidate, err := json.MarshalIndent(matchs, "", " ")
		require.NoError(t, err)
		log.Println(string(jsonCandidate))
	})
	t.Run("zero matchs with valid user", func(t *testing.T) {
		user := createNewAccount(t)
		convs, err := matchRepo.SelectMatchReqToUserId(user.ID)
		require.NoError(t, err)
		assert.Empty(t, convs)
	})
	t.Run("zero matchs with invalid user", func(t *testing.T) {
		convs, err := matchRepo.SelectMatchReqToUserId(util.RandomUUID())
		require.NoError(t, err)
		assert.Empty(t, convs)
	})
}
func Test_GetMatchById(t *testing.T) {
	matchRepo := repository.NewMatch(testQuery)
	t.Run("valid select", func(t *testing.T) {
		matchId := createNewMatch(t)
		matchRes, err := matchRepo.GetMatchById(matchId)
		require.NoError(t, err)
		assert.Equal(t, matchId, matchRes.Id)
	})
	t.Run("invalid userId", func(t *testing.T) {
		matchRes, err := matchRepo.GetMatchById(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrResourceNotFound)
		assert.Empty(t, matchRes)
	})
}

func Test_UpdateMatchById(t *testing.T) {
	matchRepo := repository.NewMatch(testQuery)
	setupFunc := func(t *testing.T) entity.Match {
		fromUsr := createNewAccount(t)
		toUsr := createNewAccount(t)
		matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, domain.Accepted)
		require.NoError(t, err)
		require.NotEmpty(t, matchId)
		match, err := matchRepo.GetMatchById(matchId)
		require.NoError(t, err)
		return match
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
		assert.ErrorIs(t, err, common.ErrResourceNotFound)
	})

}
func createNewMatch(t *testing.T) string {
	matchRepo := repository.NewMatch(testQuery)
	fromUsr := createNewAccount(t)
	toUsr := createNewAccount(t)
	matchStatus := domain.Unknown
	if util.RandomBool() {
		matchStatus = domain.Accepted
	} else {
		matchStatus = domain.Declined
	}
	matchId, err := matchRepo.InsertNewMatch(fromUsr.ID, toUsr.ID, matchStatus)
	require.NoError(t, err)
	return matchId
}
