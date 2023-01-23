package repository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apiError "github.com/xyedo/blindate/pkg/common/error"
	"github.com/xyedo/blindate/pkg/common/util"
	"github.com/xyedo/blindate/pkg/infra/repository"
)

var (
	jwtSecret = "jwt-test"
)

func Test_AddRefreshToken(t *testing.T) {
	createNewAccount(t)
}

func Test_VerifyRefreshToken(t *testing.T) {
	auth := repository.NewAuth(testQuery)
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
		assert.ErrorIs(t, err, apiError.ErrNotMatchCredential)
	})
}

func Test_DeleteRefreshToken(t *testing.T) {
	auth := repository.NewAuth(testQuery)
	t.Run("Valid Token", func(t *testing.T) {
		token := createNewToken(t)
		err := auth.DeleteRefreshToken(token)
		assert.NoError(t, err)
	})
	t.Run("Invalid token", func(t *testing.T) {
		token, err := util.RandomToken(jwtSecret, 10*time.Second)
		assert.NoError(t, err)
		err = auth.DeleteRefreshToken(token)
		require.Error(t, err)
		assert.ErrorIs(t, err, apiError.ErrResourceNotFound)
	})
}

func createNewToken(t *testing.T) string {
	auth := repository.NewAuth(testQuery)
	token, err := util.RandomToken(jwtSecret, 15*time.Second)
	assert.NoError(t, err)

	err = auth.AddRefreshToken(token)
	assert.NoError(t, err)
	return token
}
