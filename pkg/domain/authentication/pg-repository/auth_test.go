package pgrepository_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/internal/testutil"
	apperror "github.com/xyedo/blindate/pkg/common/app-error"
	authRepo "github.com/xyedo/blindate/pkg/domain/authentication/pg-repository"

	"github.com/xyedo/blindate/pkg/infrastructure"
	"github.com/xyedo/blindate/pkg/infrastructure/postgre"
)

var (
	cfg       infrastructure.Config
	jwtSecret = "jwt-test"
)

func init() {
	cfg.LoadConfig(".env.dev")
}

func Test_AddRefreshToken(t *testing.T) {
	t.Run("valid_refresh_token", func(t *testing.T) {
		db, err := postgre.OpenDB(cfg)
		require.NoError(t, err)
		authRepo := authRepo.New(db)
		testutil.CreateNewToken(t, authRepo, jwtSecret)
	})

}

func Test_VerifyRefreshToken(t *testing.T) {
	db, err := postgre.OpenDB(cfg)
	require.NoError(t, err)
	auth := authRepo.New(db)

	t.Run("Valid Token", func(t *testing.T) {
		token := testutil.CreateNewToken(t, auth, jwtSecret)

		err := auth.VerifyRefreshToken(token)
		assert.NoError(t, err)
	})
	t.Run("Invalid token", func(t *testing.T) {
		token, err := testutil.RandomToken(jwtSecret, 10*time.Second)
		assert.NoError(t, err)
		err = auth.VerifyRefreshToken(token)
		assert.Error(t, err)
		assert.ErrorIs(t, err, apperror.ErrUnauthorized)
	})
}

func Test_DeleteRefreshToken(t *testing.T) {
	db, err := postgre.OpenDB(cfg)
	require.NoError(t, err)
	auth := authRepo.New(db)
	t.Run("Valid Token", func(t *testing.T) {
		token := testutil.CreateNewToken(t, auth, jwtSecret)
		err := auth.DeleteRefreshToken(token)
		assert.NoError(t, err)
	})
	t.Run("Invalid token", func(t *testing.T) {
		token, err := testutil.RandomToken(jwtSecret, 10*time.Second)
		assert.NoError(t, err)
		err = auth.DeleteRefreshToken(token)
		require.Error(t, err)
		assert.ErrorIs(t, err, apperror.ErrNotFound)
	})
}
