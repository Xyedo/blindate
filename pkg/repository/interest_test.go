package repository

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertInterest(t *testing.T) {
	t.Run("Valid UserId", func(t *testing.T) {

		createNewInterest(t)
	})
	t.Run("Invalid UserId", func(t *testing.T) {
		repo := NewInterest(testQuery)
		intr := createNewInterest(t)
		intr.UserId = util.RandomUUID()
		err := repo.InsertInterest(intr)
		assert.Error(t, err)
		var pqErr *pq.Error
		assert.ErrorAs(t, err, &pqErr)
		assert.Equal(t, pq.ErrorCode("23503"), pqErr.Code)
		assert.True(t, strings.Contains(pqErr.Constraint, "user_id"))
	})
}

func Test_GetInterest(t *testing.T) {
	repo := NewInterest(testQuery)
	t.Run("Valid Id", func(t *testing.T) {
		expect := createNewInterest(t)
		res, err := repo.GetInterest(expect.UserId)
		assert.NoError(t, err)
		assert.Equal(t, expect.Id, res.Id)
		assert.Equal(t, expect.UserId, res.UserId)
		assert.Equal(t, expect.Bio, res.Bio)
		if len(res.Hobbies) > 0 {
			assert.Equal(t, expect.Hobbies[0].Hobbie, res.Hobbies[0].Hobbie)
			assert.NotZero(t, res.Hobbies[0].Id)
		}
		if len(res.MovieSeries) > 0 {
			assert.Equal(t, expect.MovieSeries[0].MovieSerie, res.MovieSeries[0].MovieSerie)
			assert.NotZero(t, res.MovieSeries[0].Id)
		}
		if len(res.Travels) > 0 {
			assert.Equal(t, expect.Travels[0].Travel, res.Travels[0].Travel)
			assert.NotZero(t, res.Travels[0].Id)
		}

		if len(res.Sports) > 0 {
			assert.Equal(t, expect.Sports[0].Sport, res.Sports[0].Sport)
			assert.NotZero(t, res.Sports[0].Id)
		}
	})
	t.Run("Invalid Id", func(t *testing.T) {
		res, err := repo.GetInterest(util.RandomUUID())
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Nil(t, res)
	})
}
func Test_UpdateInterest(t *testing.T) {
	repo := NewInterest(testQuery)
	intr := createNewInterest(t)
	intr, err := repo.GetInterest(intr.UserId)
	assert.NoError(t, err)
	err = repo.UpdateInterest(intr)
	assert.NoError(t, err)
}

func createNewInterest(t *testing.T) *domain.Bio {
	repo := NewInterest(testQuery)
	user := createNewAccount(t)
	hobbie := make([]domain.Hobbie, 0)
	for i := 0; i < int(util.RandomInt(0, 10)); i++ {
		hobbie = append(hobbie, domain.Hobbie{
			Hobbie: util.RandomString(12),
		})
	}
	moviesSeries := make([]domain.MovieSerie, 0)
	for i := 0; i < int(util.RandomInt(0, 10)); i++ {
		moviesSeries = append(moviesSeries, domain.MovieSerie{
			MovieSerie: util.RandomString(12),
		})
	}
	traveling := make([]domain.Travel, 0)
	for i := 0; i < int(util.RandomInt(0, 10)); i++ {
		traveling = append(traveling, domain.Travel{
			Travel: util.RandomString(12),
		})
	}
	sports := make([]domain.Sport, 0)
	for i := 0; i < int(util.RandomInt(0, 10)); i++ {
		sports = append(sports, domain.Sport{
			Sport: util.RandomString(12),
		})
	}
	interest := &domain.Bio{
		UserId:      user.ID,
		Hobbies:     hobbie,
		MovieSeries: moviesSeries,
		Travels:     traveling,
		Sports:      sports,
		Bio:         util.RandomString(50),
	}
	err := repo.InsertInterest(interest)
	assert.NoError(t, err)
	assert.NotNil(t, interest.Id)
	return interest
}
