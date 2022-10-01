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
		createNewLocation(t)
	})
	t.Run("Invalid user_id", func(t *testing.T) {
		location := createNewLocation(t)
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
	location := createNewLocation(t)
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
	expected := createNewLocation(t)
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
	fromUser := createNewLocation(t)

	for i := 0; i < limit*2; i++ {
		createNewLocation(t)
	}
	users, err := repo.GetClosestUser(fromUser.Geog, limit)
	assert.NoError(t, err)
	assert.NotZero(t, users)
	assert.Len(t, users, limit)
	log.Println(users)
}

func createNewLocation(t *testing.T) *entity.Location {
	repo := NewLocation(testQuery)
	user := createNewAccount(t)
	location := entity.Location{
		UserId: user.ID,
		Geog:   util.RandomPoint(5),
	}
	row, err := repo.InsertNewLocation(&location)
	assert.NoError(t, err)
	assert.Equal(t, 1, int(row))
	return &location
}
