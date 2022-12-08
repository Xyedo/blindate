package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xyedo/blindate/pkg/common"
	"github.com/xyedo/blindate/pkg/domain"
	"github.com/xyedo/blindate/pkg/repository"
	"github.com/xyedo/blindate/pkg/util"
)

func Test_InsertBio(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid UserId", func(t *testing.T) {
		user := createNewAccount(t)
		createNewInterestBio(t, user.ID)
	})
	t.Run("Invalid UserId", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		bio.UserId = util.RandomUUID()
		err := repo.InsertInterestBio(bio)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("invalid on unique constraint", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		err := repo.InsertInterestBio(bio)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)
	})
}
func Test_SelectBio(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid UserId", func(t *testing.T) {
		user := createNewAccount(t)
		exp := createNewInterestBio(t, user.ID)
		res, err := repo.SelectInterestBio(user.ID)
		require.NoError(t, err)
		assert.Equal(t, exp.Id, res.Id)
		assert.Equal(t, exp.UserId, res.UserId)
		assert.Equal(t, exp.Bio, res.Bio)
	})
	t.Run("Invalid UserId", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		bio.UserId = util.RandomUUID()
		err := repo.InsertInterestBio(bio)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
}
func Test_UpdateInterestBio(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid", func(t *testing.T) {

		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		bio.Bio = util.RandomString(12)
		err := repo.UpdateInterestBio(*bio)
		require.NoError(t, err)
	})
	t.Run("Ivalid userId", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		bio.UserId = util.RandomUUID()
		bio.Bio = util.RandomString(12)
		err := repo.UpdateInterestBio(*bio)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrResourceNotFound)
	})
}

func Test_InsertHobbies(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid InterestId", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		createNewInterestHobbie(t, bio.Id)

	})
	t.Run("Invalid InterestId", func(t *testing.T) {
		hobbies := make([]domain.Hobbie, 0)
		for i := 0; i < int(util.RandomInt(1, 10)); i++ {
			hobbies = append(hobbies, domain.Hobbie{
				Hobbie: util.RandomString(12),
			})
		}
		err := repo.InsertInterestHobbies(util.RandomUUID(), hobbies)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("Invalid Unique", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		hobbie := createNewInterestHobbie(t, intr.Id)
		err := repo.InsertInterestHobbies(intr.Id, hobbie)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)
	})
	t.Run("Too much bio", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		err := repo.InsertNewStats(bio.Id)
		require.NoError(t, err)
		hobbies := make([]domain.Hobbie, 0)
		for i := 0; i < 11; i++ {
			hobbies = append(hobbies, domain.Hobbie{
				Hobbie: util.RandomString(12),
			})
		}
		err = repo.InsertInterestHobbies(bio.Id, hobbies)
		require.Error(t, err)
		require.Implements(t, (*common.APIError)(nil), err)
	})
}
func Test_UpdateHobbies(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid Update", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		hobbies := createNewInterestHobbie(t, intr.Id)
		for i := range hobbies {
			hobbies[i].Hobbie = util.RandomString(7)
		}
		for i := 0; i < 10-len(hobbies); i++ {
			hobbies = append(hobbies, domain.Hobbie{
				Hobbie: util.RandomString(15),
			})
		}
		err := repo.UpdateInterestHobbies(intr.Id, hobbies)
		assert.NoError(t, err)
	})
	t.Run("Valid Update but over than 10", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		err := repo.InsertNewStats(intr.Id)
		require.NoError(t, err)
		hobbies := createNewInterestHobbie(t, intr.Id)
		for i := range hobbies {
			hobbies[i].Hobbie = util.RandomString(7)
		}
		for i := 0; i < 10; i++ {
			hobbies = append(hobbies, domain.Hobbie{
				Hobbie: util.RandomString(15),
			})
		}
		err = repo.UpdateInterestHobbies(intr.Id, hobbies)
		require.Error(t, err)
		require.Implements(t, (*common.APIError)(nil), err)
	})
}
func Test_DeleteHobbies(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid Update", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		hobbies := createNewInterestHobbie(t, intr.Id)
		ids := make([]string, 0, len(hobbies))
		for _, hobie := range hobbies {
			ids = append(ids, hobie.Id)
		}
		err := repo.DeleteInterestHobbies(intr.Id, ids)
		require.NoError(t, err)

	})
}
func Test_InsertMovieSeries(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid InterestId", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		createNewInterestMovieSeries(t, bio.Id)

	})
	t.Run("Invalid InterestId", func(t *testing.T) {
		movieSeries := make([]domain.MovieSerie, 0)
		for i := 0; i < int(util.RandomInt(1, 10)); i++ {
			movieSeries = append(movieSeries, domain.MovieSerie{
				MovieSerie: util.RandomString(12),
			})
		}
		err := repo.InsertInterestMovieSeries(util.RandomUUID(), movieSeries)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("Invalid Unique", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		movieSeries := createNewInterestMovieSeries(t, intr.Id)
		err := repo.InsertInterestMovieSeries(intr.Id, movieSeries)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)
	})
	t.Run("Too much movies", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		err := repo.InsertNewStats(bio.Id)
		require.NoError(t, err)
		movieSeries := make([]domain.MovieSerie, 0)
		for i := 0; i <= 11; i++ {
			movieSeries = append(movieSeries, domain.MovieSerie{
				MovieSerie: util.RandomString(12),
			})
		}
		err = repo.InsertInterestMovieSeries(bio.Id, movieSeries)
		require.Error(t, err)
		require.Implements(t, (*common.APIError)(nil), err)
	})
}
func Test_UpdateMovieSeries(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid Update", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		movieSeries := createNewInterestMovieSeries(t, intr.Id)
		for i := range movieSeries {
			movieSeries[i].MovieSerie = util.RandomString(19)
		}
		for i := 0; i < 10-len(movieSeries); i++ {
			movieSeries = append(movieSeries, domain.MovieSerie{
				MovieSerie: util.RandomString(15),
			})
		}
		err := repo.UpdateInterestMovieSeries(intr.Id, movieSeries)
		require.NoError(t, err)
	})
	t.Run("Valid Update but over than 10", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		err := repo.InsertNewStats(intr.Id)
		require.NoError(t, err)
		movieSeries := createNewInterestMovieSeries(t, intr.Id)
		for i := range movieSeries {
			movieSeries[i].MovieSerie = util.RandomString(7)
		}
		for i := 0; i < 11; i++ {
			movieSeries = append(movieSeries, domain.MovieSerie{
				MovieSerie: util.RandomString(15),
			})
		}
		err = repo.UpdateInterestMovieSeries(intr.Id, movieSeries)
		require.Error(t, err)
		require.Implements(t, (*common.APIError)(nil), err)
	})
}
func Test_DeleteMovieSeries(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid Update", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		movieSeries := createNewInterestMovieSeries(t, intr.Id)
		ids := make([]string, 0, len(movieSeries))
		for _, hobie := range movieSeries {
			ids = append(ids, hobie.Id)
		}
		err := repo.DeleteInterestMovieSeries(intr.Id, ids)
		require.NoError(t, err)
	})
}

func Test_InsertTraveling(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid InterestId", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		createNewInterestTraveling(t, bio.Id)

	})
	t.Run("Invalid InterestId", func(t *testing.T) {
		travels := make([]domain.Travel, 0)
		for i := 0; i < int(util.RandomInt(1, 10)); i++ {
			travels = append(travels, domain.Travel{
				Travel: util.RandomString(12),
			})
		}
		err := repo.InsertInterestTraveling(util.RandomUUID(), travels)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("Invalid Unique", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		travels := createNewInterestTraveling(t, intr.Id)
		err := repo.InsertInterestTraveling(intr.Id, travels)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)
	})
	t.Run("Too much travels", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		err := repo.InsertNewStats(bio.Id)
		require.NoError(t, err)
		travels := make([]domain.Travel, 0)
		for i := 0; i < 11; i++ {
			travels = append(travels, domain.Travel{
				Travel: util.RandomString(12),
			})
		}
		err = repo.InsertInterestTraveling(bio.Id, travels)
		require.Error(t, err)
		require.Implements(t, (*common.APIError)(nil), err)
	})
}
func Test_UpdateTraveling(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid Update", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		travels := createNewInterestTraveling(t, intr.Id)
		for i := range travels {
			travels[i].Travel = util.RandomString(19)
		}
		for i := 0; i < 10-len(travels); i++ {
			travels = append(travels, domain.Travel{
				Travel: util.RandomString(15),
			})
		}
		err := repo.UpdateInterestTraveling(intr.Id, travels)
		assert.NoError(t, err)
	})
	t.Run("Valid Update but over than 10", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		err := repo.InsertNewStats(intr.Id)
		require.NoError(t, err)
		travels := createNewInterestTraveling(t, intr.Id)
		for i := range travels {
			travels[i].Travel = util.RandomString(7)
		}
		for i := 0; i < 10; i++ {
			travels = append(travels, domain.Travel{
				Travel: util.RandomString(15),
			})
		}
		err = repo.UpdateInterestTraveling(intr.Id, travels)
		require.Error(t, err)
		require.Implements(t, (*common.APIError)(nil), err)
	})
}
func Test_DeleteTraveling(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid Update", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		travels := createNewInterestTraveling(t, intr.Id)
		ids := make([]string, 0, len(travels))
		for _, hobie := range travels {
			ids = append(ids, hobie.Id)
		}
		err := repo.DeleteInterestTraveling(intr.Id, ids)
		assert.NoError(t, err)
	})
}
func Test_InsertSports(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid InterestId", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		createNewInterestSport(t, bio.Id)

	})
	t.Run("Invalid InterestId", func(t *testing.T) {
		sports := make([]domain.Sport, 0)
		for i := 0; i < int(util.RandomInt(1, 10)); i++ {
			sports = append(sports, domain.Sport{
				Sport: util.RandomString(12),
			})
		}
		err := repo.InsertInterestSports(util.RandomUUID(), sports)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrRefNotFound23503)
	})
	t.Run("Invalid Unique", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		sports := createNewInterestSport(t, intr.Id)
		err := repo.InsertInterestSports(intr.Id, sports)
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrUniqueConstraint23505)
	})
	t.Run("Too much sports", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		err := repo.InsertNewStats(bio.Id)
		require.NoError(t, err)
		sports := make([]domain.Sport, 0)
		for i := 0; i < 11; i++ {
			sports = append(sports, domain.Sport{
				Sport: util.RandomString(12),
			})
		}
		err = repo.InsertInterestSports(bio.Id, sports)
		require.Error(t, err)
		require.Implements(t, (*common.APIError)(nil), err)
	})
}
func Test_UpdateSports(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid Update", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		sports := createNewInterestSport(t, intr.Id)
		for i := range sports {
			sports[i].Sport = util.RandomString(19)
		}
		for i := 0; i < 10-len(sports); i++ {
			sports = append(sports, domain.Sport{
				Sport: util.RandomString(15),
			})
		}
		err := repo.UpdateInterestSport(intr.Id, sports)
		require.NoError(t, err)
	})
	t.Run("Valid Update but over than 10", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		err := repo.InsertNewStats(intr.Id)
		require.NoError(t, err)
		sports := createNewInterestSport(t, intr.Id)
		for i := range sports {
			sports[i].Sport = util.RandomString(7)
		}
		for i := 0; i < 10; i++ {
			sports = append(sports, domain.Sport{
				Sport: util.RandomString(15),
			})
		}
		err = repo.UpdateInterestSport(intr.Id, sports)
		require.Error(t, err)
		require.Implements(t, (*common.APIError)(nil), err)
	})
}
func Test_DeleteSports(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid Update", func(t *testing.T) {
		user := createNewAccount(t)
		intr := createNewInterestBio(t, user.ID)
		sports := createNewInterestSport(t, intr.Id)
		ids := make([]string, 0, len(sports))
		for _, hobie := range sports {
			ids = append(ids, hobie.Id)
		}
		err := repo.DeleteInterestSports(intr.Id, ids)
		assert.NoError(t, err)
	})
}

func Test_GetInterest(t *testing.T) {
	repo := repository.NewInterest(testQuery)
	t.Run("Valid Id AND Using All", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		hobbies := createNewInterestHobbie(t, bio.Id)
		movieSeries := createNewInterestMovieSeries(t, bio.Id)
		travels := createNewInterestTraveling(t, bio.Id)
		sports := createNewInterestSport(t, bio.Id)

		res, err := repo.GetInterest(bio.UserId)
		require.NoError(t, err)
		assert.Equal(t, bio.Id, res.Id)
		assert.Equal(t, bio.UserId, res.UserId)
		assert.Equal(t, bio.Bio, res.Bio.Bio)
		if len(res.Hobbies) > 0 {
			assert.NotZero(t, hobbies[0].Hobbie)
			assert.NotZero(t, res.Hobbies[0].Hobbie)
			assert.NotZero(t, res.Hobbies[0].Id)
		}
		if len(res.MovieSeries) > 0 {
			assert.NotZero(t, movieSeries[0].MovieSerie)
			assert.NotZero(t, res.MovieSeries[0].MovieSerie)
			assert.NotZero(t, res.MovieSeries[0].Id)
		}
		if len(res.Travels) > 0 {
			assert.NotZero(t, travels[0].Travel)
			assert.NotZero(t, res.Travels[0].Travel)
			assert.NotZero(t, res.Travels[0].Id)
		}

		if len(res.Sports) > 0 {
			assert.NotZero(t, sports[0].Sport)
			assert.NotZero(t, res.Sports[0].Sport)
			assert.NotZero(t, res.Sports[0].Id)
		}
	})
	t.Run("Valid Id But Partial hobbies", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		hobbies := createNewInterestHobbie(t, bio.Id)
		res, err := repo.GetInterest(bio.UserId)
		require.NoError(t, err)
		assert.Equal(t, bio.Id, res.Id)
		assert.Equal(t, bio.UserId, res.UserId)
		assert.Equal(t, bio.Bio, res.Bio.Bio)

		assert.NotZero(t, hobbies[0].Hobbie)
		assert.NotZero(t, res.Hobbies[0].Hobbie)
		assert.NotZero(t, res.Hobbies[0].Id)

		assert.Zero(t, res.MovieSeries)
		assert.Zero(t, res.Travels)
		assert.Zero(t, res.Sports)

	})
	t.Run("Valid Id But Partial MovieSeries", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		movieSeries := createNewInterestMovieSeries(t, bio.Id)
		res, err := repo.GetInterest(bio.UserId)
		require.NoError(t, err)
		assert.Equal(t, bio.Id, res.Id)
		assert.Equal(t, bio.UserId, res.UserId)
		assert.Equal(t, bio.Bio, res.Bio.Bio)

		assert.NotZero(t, movieSeries[0].MovieSerie)
		assert.NotZero(t, res.MovieSeries[0].MovieSerie)
		assert.NotZero(t, res.MovieSeries[0].Id)

		assert.Zero(t, res.Hobbies)
		assert.Zero(t, res.Travels)
		assert.Zero(t, res.Sports)

	})
	t.Run("Valid Id But Partial travels", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		travels := createNewInterestTraveling(t, bio.Id)
		res, err := repo.GetInterest(bio.UserId)
		require.NoError(t, err)
		assert.Equal(t, bio.Id, res.Id)
		assert.Equal(t, bio.UserId, res.UserId)
		assert.Equal(t, bio.Bio, res.Bio.Bio)

		assert.NotZero(t, travels[0].Travel)
		assert.NotZero(t, res.Travels[0].Travel)
		assert.NotZero(t, res.Travels[0].Id)

		assert.Zero(t, res.Hobbies)
		assert.Zero(t, res.MovieSeries)
		assert.Zero(t, res.Sports)

	})
	t.Run("Valid Id But Partial sports", func(t *testing.T) {
		user := createNewAccount(t)
		bio := createNewInterestBio(t, user.ID)
		sports := createNewInterestSport(t, bio.Id)
		res, err := repo.GetInterest(bio.UserId)
		require.NoError(t, err)
		assert.Equal(t, bio.Id, res.Id)
		assert.Equal(t, bio.UserId, res.UserId)
		assert.Equal(t, bio.Bio, res.Bio.Bio)

		assert.NotZero(t, sports[0].Sport)
		assert.NotZero(t, res.Sports[0].Sport)
		assert.NotZero(t, res.Sports[0].Id)

		assert.Zero(t, res.Hobbies)
		assert.Zero(t, res.MovieSeries)
		assert.Zero(t, res.Travels)

	})
	t.Run("Invalid Id", func(t *testing.T) {
		res, err := repo.GetInterest(util.RandomUUID())
		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrResourceNotFound)
		assert.Zero(t, res)
	})
}

func createNewInterestBio(t *testing.T, userId string) *domain.Bio {
	repo := repository.NewInterest(testQuery)
	bio := &domain.Bio{
		UserId: userId,
		Bio:    util.RandomString(50),
	}
	err := repo.InsertInterestBio(bio)
	assert.NoError(t, err)
	assert.NotNil(t, bio.Id)
	return bio
}

func createNewInterestHobbie(t *testing.T, interestId string) []domain.Hobbie {
	repo := repository.NewInterest(testQuery)
	hobbies := make([]domain.Hobbie, 0)
	for i := 0; i < int(util.RandomInt(1, 10)); i++ {
		hobbies = append(hobbies, domain.Hobbie{
			Hobbie: util.RandomString(12),
		})
	}

	err := repo.InsertInterestHobbies(interestId, hobbies)
	assert.NoError(t, err)
	assert.NotZero(t, hobbies[0].Id)
	return hobbies
}
func createNewInterestMovieSeries(t *testing.T, interestId string) []domain.MovieSerie {
	repo := repository.NewInterest(testQuery)
	moviesSeries := make([]domain.MovieSerie, 0)
	for i := 0; i < int(util.RandomInt(1, 10)); i++ {
		moviesSeries = append(moviesSeries, domain.MovieSerie{
			MovieSerie: util.RandomString(12),
		})
	}

	err := repo.InsertInterestMovieSeries(interestId, moviesSeries)
	assert.NoError(t, err)
	return moviesSeries
}
func createNewInterestTraveling(t *testing.T, interestId string) []domain.Travel {
	repo := repository.NewInterest(testQuery)
	travels := make([]domain.Travel, 0)
	for i := 0; i < int(util.RandomInt(1, 10)); i++ {
		travels = append(travels, domain.Travel{
			Travel: util.RandomString(12),
		})
	}

	err := repo.InsertInterestTraveling(interestId, travels)
	assert.NoError(t, err)
	return travels
}
func createNewInterestSport(t *testing.T, interestId string) []domain.Sport {
	repo := repository.NewInterest(testQuery)
	sports := make([]domain.Sport, 0)
	for i := 0; i < int(util.RandomInt(1, 10)); i++ {
		sports = append(sports, domain.Sport{
			Sport: util.RandomString(12),
		})
	}

	err := repo.InsertInterestSports(interestId, sports)
	assert.NoError(t, err)
	return sports
}
