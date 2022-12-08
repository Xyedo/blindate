package repository_test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain/entity"
	"github.com/xyedo/blindate/pkg/repository"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertNewLocation(t *testing.T) {
	repo := repository.NewLocation(testQuery)
	t.Run("Valid Create new Location", func(t *testing.T) {
		user := createNewAccount(t)
		createNewLocation(t, user.ID)
	})
	t.Run("Valid but Double", func(t *testing.T) {
		user := createNewAccount(t)
		location := createNewLocation(t, user.ID)
		err := repo.InsertNewLocation(location)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)

	})
	t.Run("Invalid user_id", func(t *testing.T) {
		user := createNewAccount(t)
		location := createNewLocation(t, user.ID)
		location.UserId = "e590666c-3ea8-4fda-958c-c2dc6c2599b6"
		err := repo.InsertNewLocation(location)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
}

func Test_UpdateLocation(t *testing.T) {
	repo := repository.NewLocation(testQuery)
	user := createNewAccount(t)
	location := createNewLocation(t, user.ID)
	t.Run("Valid Update Location", func(t *testing.T) {

		location.Geog = util.RandomPoint(5)
		err := repo.UpdateLocation(location)
		assert.NoError(t, err)
	})
	t.Run("Invalid UserId", func(t *testing.T) {
		location.Geog = util.RandomPoint(5)
		location.UserId = "e590666c-3ea8-4fda-958c-c2dc6c2599b6"
		err := repo.UpdateLocation(location)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrResourceNotFound)

	})
}

func Test_GetLocationByUserId(t *testing.T) {
	repo := repository.NewLocation(testQuery)
	user := createNewAccount(t)
	expected := createNewLocation(t, user.ID)
	t.Run("valid Getter", func(t *testing.T) {
		_, err := repo.GetLocationByUserId(expected.UserId)
		assert.NoError(t, err)
	})

	t.Run("Invalid User_id", func(t *testing.T) {
		location, err := repo.GetLocationByUserId("e590666c-3ea8-4fda-958c-c2dc6c2599b6")
		require.Error(t, err)
		require.ErrorIs(t, err, common.ErrResourceNotFound)
		assert.Zero(t, location)
	})

}

func Test_GetClosestUser(t *testing.T) {
	limit := 5
	repo := repository.NewLocation(testQuery)
	user := createNewAccount(t)
	fromUser := createNewLocation(t, user.ID)

	for i := 0; i < limit*2; i++ {
		useri := createNewAccount(t)
		createNewLocation(t, useri.ID)
	}
	log.Println("request User", user.ID)
	candidateMatch, err := repo.GetClosestUser(user.ID, fromUser.Geog, limit)
	require.NoError(t, err)
	assert.NotZero(t, candidateMatch)
	assert.Len(t, candidateMatch, limit)
	jsonCandidate, err := json.Marshal(candidateMatch)
	require.NoError(t, err)
	log.Println(string(jsonCandidate))

}

func createNewLocation(t *testing.T, userId string) *entity.Location {
	repo := repository.NewLocation(testQuery)

	location := entity.Location{
		UserId: userId,
		Geog:   util.RandomPoint(5),
	}
	err := repo.InsertNewLocation(&location)
	require.NoError(t, err)
	return &location
}
