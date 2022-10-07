package repository

import (
	"database/sql"
	"log"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/entity"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertNewLocation(t *testing.T) {
	repo := NewLocation(testQuery)
	t.Run("Valid Create new Location", func(t *testing.T) {
		user := createNewAccount(t)
		createNewLocation(t, user.ID)
	})
	t.Run("Valid but Double", func(t *testing.T) {
		user := createNewAccount(t)
		location := createNewLocation(t, user.ID)
		row, err := repo.InsertNewLocation(location)
		assert.Error(t, err)
		var pqErr *pq.Error
		if assert.ErrorAs(t, err, &pqErr) {
			assert.Equal(t, pq.ErrorCode("23505"), pqErr.Code)
			assert.Contains(t, pqErr.Constraint, "locations_pkey")
		}
		assert.Zero(t, row)

	})
	t.Run("Invalid user_id", func(t *testing.T) {
		user := createNewAccount(t)
		location := createNewLocation(t, user.ID)
		location.UserId = "e590666c-3ea8-4fda-958c-c2dc6c2599b6"
		row, err := repo.InsertNewLocation(location)
		var pqError *pq.Error

		assert.ErrorAs(t, err, &pqError)
		assert.Equal(t, pq.ErrorCode("23503"), pqError.Code)
		assert.True(t, strings.Contains(pqError.Constraint, "user_id"))
		assert.Zero(t, row)

	})
}

func Test_UpdateLocation(t *testing.T) {
	repo := NewLocation(testQuery)
	user := createNewAccount(t)
	location := createNewLocation(t, user.ID)
	t.Run("Valid Update Location", func(t *testing.T) {

		location.Geog = util.RandomPoint(5)
		row, err := repo.UpdateLocation(location)
		assert.NoError(t, err)
		assert.Equal(t, 1, int(row))
	})
	t.Run("Invalid UserId", func(t *testing.T) {
		location.Geog = util.RandomPoint(5)
		location.UserId = "e590666c-3ea8-4fda-958c-c2dc6c2599b6"
		row, err := repo.UpdateLocation(location)
		assert.NoError(t, err)
		assert.Zero(t, row)

	})
}

func Test_GetLocationByUserId(t *testing.T) {
	repo := NewLocation(testQuery)
	user := createNewAccount(t)
	expected := createNewLocation(t, user.ID)
	t.Run("valid Getter", func(t *testing.T) {
		_, err := repo.GetLocationByUserId(expected.UserId)
		assert.NoError(t, err)
	})

	t.Run("Invalid User_id", func(t *testing.T) {
		location, err := repo.GetLocationByUserId("e590666c-3ea8-4fda-958c-c2dc6c2599b6")
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Zero(t, location)
	})

}

func Test_GetClosestUser(t *testing.T) {
	limit := 5
	repo := NewLocation(testQuery)
	user := createNewAccount(t)
	fromUser := createNewLocation(t, user.ID)

	for i := 0; i < limit*2; i++ {
		useri := createNewAccount(t)
		createNewLocation(t, useri.ID)
	}
	users, err := repo.GetClosestUser(fromUser.Geog, limit)
	assert.NoError(t, err)
	assert.NotZero(t, users)
	assert.Len(t, users, limit)
	log.Println(users)
}

func createNewLocation(t *testing.T, userId string) *entity.Location {
	repo := NewLocation(testQuery)

	location := entity.Location{
		UserId: userId,
		Geog:   util.RandomPoint(5),
	}
	row, err := repo.InsertNewLocation(&location)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(row))
	return &location
}
