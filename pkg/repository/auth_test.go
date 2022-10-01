package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/util"
)

var (
	jwtSecret = "jwt-test"
)

func Test_AddRefreshToken(t *testing.T) {
	createNewAccount(t)
}

func Test_VerifyRefreshToken(t *testing.T) {
	auth := NewAuth(testQuery)
	t.Run("Valid Token", func(t *testing.T) {
		token := createNewToken(t)

		err := auth.VerifyRefreshToken(token)
		assert.NoError(t, err)
	})
	t.Run("Invalid token", func(t *testing.T) {
		token, err := util.RandomToken(jwtSecret, 10*time.Second)
		assert.NoError(t, err)
		err = auth.VerifyRefreshToken(token)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}

func Test_DeleteRefreshToken(t *testing.T) {
	auth := NewAuth(testQuery)
	t.Run("Valid Token", func(t *testing.T) {
		token := createNewToken(t)
		rows, err := auth.DeleteRefreshToken(token)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), rows)
	})
	t.Run("Invalid token", func(t *testing.T) {
		token, err := util.RandomToken(jwtSecret, 10*time.Second)
		assert.NoError(t, err)
		rows, err := auth.DeleteRefreshToken(token)
		assert.NoError(t, err)
		assert.Zero(t, rows)
	})
}

func createNewToken(t *testing.T) string {
	auth := NewAuth(testQuery)
	token, err := util.RandomToken(jwtSecret, 15*time.Second)
	assert.NoError(t, err)

	rows, err := auth.AddRefreshToken(token)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rows)
	return token
}
