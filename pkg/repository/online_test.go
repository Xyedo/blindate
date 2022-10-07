package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertNewOnline(t *testing.T) {
	repo := NewOnline(testQuery)
	t.Run("Valid And Unique", func(t *testing.T) {
		user := createNewAccount(t)
		online := createNewOnline(t, user.ID)
		err := repo.InsertNewOnline(online)
		assert.Error(t, err)
		var pqErr *pq.Error
		if assert.ErrorAs(t, err, &pqErr) {
			assert.Equal(t, pq.ErrorCode("23505"), pqErr.Code)
			assert.Contains(t, pqErr.Constraint, "onlines_pkey")
		}
	})
	t.Run("Invalid user_id", func(t *testing.T) {
		online := &domain.Online{
			UserId:     util.RandomUUID(),
			LastOnline: time.Now(),
			IsOnline:   false,
		}
		err := repo.InsertNewOnline(online)
		assert.Error(t, err)
		var pqErr *pq.Error
		if assert.ErrorAs(t, err, &pqErr) {
			assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
			assert.Contains(t, pqErr.Constraint, "user_id")
		}
	})

}

func Test_SelectOnline(t *testing.T) {
	repo := NewOnline(testQuery)
	user := createNewAccount(t)
	exp := createNewOnline(t, user.ID)
	res, err := repo.SelectOnline(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, exp.IsOnline, res.IsOnline)
	assert.Equal(t, exp.UserId, res.UserId)
	assert.NotZero(t, res.LastOnline)
}

func createNewOnline(t *testing.T, userId string) *domain.Online {
	repo := NewOnline(testQuery)
	online := &domain.Online{
		UserId:     userId,
		LastOnline: time.Now(),
		IsOnline:   false,
	}
	err := repo.InsertNewOnline(online)
	assert.NoError(t, err)
	return online
}

func Test_UpdateOnline(t *testing.T) {
	repo := NewOnline(testQuery)
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
			assert.ErrorIs(t, err, sql.ErrNoRows)
		})
		t.Run("Offline", func(t *testing.T) {
			err := repo.UpdateOnline(util.RandomUUID(), false)
			assert.ErrorIs(t, err, sql.ErrNoRows)
		})
	})

}
