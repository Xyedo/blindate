package testutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/domain/authentication"
)

func CreateNewToken(t *testing.T, auth authentication.Repository, jwtSecret string) string {
	token, err := RandomToken(jwtSecret, 15*time.Second)
	assert.NoError(t, err)

	err = auth.AddRefreshToken(token)
	assert.NoError(t, err)
	return token
}
