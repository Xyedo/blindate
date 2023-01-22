package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/common"
	onlineEntities "github.com/xyedo/blindate/pkg/domain/online/entities"
	"github.com/xyedo/blindate/pkg/infra/repository"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertNewOnline(t *testing.T) {
	repo := repository.NewOnline(testQuery)
	t.Run("Valid And Unique", func(t *testing.T) {
		user := createNewAccount(t)
		online := createNewOnline(t, user.ID)
		err := repo.InsertNewOnline(online)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)

	})
	t.Run("Invalid user_id", func(t *testing.T) {
		online := onlineEntities.DTO{
			UserId:     util.RandomUUID(),
			LastOnline: time.Now(),
			IsOnline:   false,
		}
		err := repo.InsertNewOnline(online)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})

}

func Test_SelectOnline(t *testing.T) {
	repo := repository.NewOnline(testQuery)
	t.Run("valid select", func(t *testing.T) {

		user := createNewAccount(t)
		exp := createNewOnline(t, user.ID)
		res, err := repo.SelectOnline(user.ID)
		require.NoError(t, err)
		assert.Equal(t, exp.IsOnline, res.IsOnline)
		assert.Equal(t, exp.UserId, res.UserId)
		assert.NotZero(t, res.LastOnline)
	})
	t.Run("invalid userId", func(t *testing.T) {
		res, err := repo.SelectOnline(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrResourceNotFound)
		require.Zero(t, res)
	})

}

func Test_UpdateOnline(t *testing.T) {
	repo := repository.NewOnline(testQuery)
	t.Run("Valid Id", func(t *testing.T) {
		t.Run("Online", func(t *testing.T) {
			user := createNewAccount(t)
			online := createNewOnline(t, user.ID)
			err := repo.UpdateOnline(online.UserId, true)
			assert.NoError(t, err)
		})
		t.Run("Offline", func(t *testing.T) {
			user := createNewAccount(t)
			online := createNewOnline(t, user.ID)
			err := repo.UpdateOnline(online.UserId, false)
			assert.NoError(t, err)
		})
	})
	t.Run("Invalid Id", func(t *testing.T) {
		t.Run("Online", func(t *testing.T) {
			err := repo.UpdateOnline(util.RandomUUID(), true)
			require.ErrorIs(t, err, common.ErrResourceNotFound)
		})
		t.Run("Offline", func(t *testing.T) {
			err := repo.UpdateOnline(util.RandomUUID(), false)
			require.ErrorIs(t, err, common.ErrResourceNotFound)
		})
	})

}
func createNewOnline(t *testing.T, userId string) onlineEntities.DTO {
	repo := repository.NewOnline(testQuery)
	online := onlineEntities.DTO{
		UserId:     userId,
		LastOnline: time.Now(),
		IsOnline:   false,
	}
	err := repo.InsertNewOnline(online)
	assert.NoError(t, err)
	return online
}
